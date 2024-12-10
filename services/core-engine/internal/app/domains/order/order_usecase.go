package order

import (
	"context"
	"fmt"
	"time"

	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"core-engine/internal/app/domains/order/model"
	"core-engine/internal/app/domains/user"
	serverError "core-engine/pkg/error"
	gormpkg "core-engine/pkg/gorm"
	"core-engine/pkg/jwt"
	"core-engine/pkg/kafka"
)

type usecase struct {
	redisLock       *redsync.Redsync
	writeDB         *gorm.DB
	kafkaProducer   kafka.Producer
	validator       *validator.Validate
	orderRepository model.Repository
	userRepository  user.Repository
}

// NewUsecase returns new order usecase.
func NewUsecase(
	redisLock *redsync.Redsync,
	writeDB *gorm.DB,
	kafkaProducer kafka.Producer,
	validator *validator.Validate,
	orderRepository model.Repository,
	userRepository user.Repository,
) *usecase {
	return &usecase{
		redisLock:       redisLock,
		writeDB:         writeDB,
		kafkaProducer:   kafkaProducer,
		validator:       validator,
		orderRepository: orderRepository,
		userRepository:  userRepository,
	}
}

func (u *usecase) ProcessOrder(ctx context.Context, orderReq model.OrderRequest) (model.Order, error) {
	tokenPayload := jwt.GetPayloadFromContext(ctx)

	// Check user detail
	userDetail, err := u.userRepository.FindUserByEmail(ctx, tokenPayload.Email)
	if err != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	// Check account status
	if !userDetail.Status {
		return model.Order{}, serverError.ErrUserBlocked(nil) // User already deactivated
	}

	// Check crypto pair detail
	cryptoPairDetail, err := u.orderRepository.GetPairDetail(ctx, orderReq.PairCode)
	if err != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	targetCryptoID := cryptoPairDetail.PrimaryCryptoID
	if orderReq.Side == model.OrderSideBuy {
		// When buying, check if user have enough secondary balance for buying primary crypto
		targetCryptoID = cryptoPairDetail.SecondaryCryptoID
	}

	// Lock all balance activity for this specific user
	mutex := u.redisLock.NewMutex(fmt.Sprintf("locking#member#%v#%v", userDetail.ID, targetCryptoID))
	if err := mutex.Lock(); err != nil {
		log.Context(ctx).Error(err)
		return model.Order{}, err
	}

	defer func() { // Release the lock so other processes or threads can obtain a lock.
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Context(ctx).Error(err)
		}
	}()

	userWallet, err := u.orderRepository.GetUserWallet(ctx, userDetail.ID, targetCryptoID)
	if err != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	// Validate user balance
	if !userWallet.IsEnoughBalance(orderReq) {
		return model.Order{}, serverError.ErrInsufficientBalance(nil)
	}

	// Deduct user wallet balance
	var errBalanceUpdate error
	switch orderReq.Side {
	case model.OrderSideSell:
		errBalanceUpdate = u.orderRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -orderReq.Quantity)

	case model.OrderSideBuy:
		totalAmount := orderReq.Price * orderReq.Quantity
		errBalanceUpdate = u.orderRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -totalAmount)
	}

	if errBalanceUpdate != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	// Save to table orders
	newOrder := model.Order{
		UserID:          userDetail.ID,
		PairID:          cryptoPairDetail.ID,
		Quantity:        orderReq.Quantity,
		Price:           orderReq.Price,
		Type:            orderReq.Type,
		Side:            orderReq.Side,
		Status:          model.OrderStatusProgress,
		TransactionTime: time.Now().Unix(),
	}

	order, err := u.orderRepository.SaveOrder(ctx, newOrder)
	if err != nil {
		return model.Order{}, err
	}

	// Publish to matching engine
	if err := u.kafkaProducer.Send(ctx, cryptoPairDetail.Code, cast.ToString(order.ID), order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (u *usecase) MatchOrder(ctx context.Context, tradeReq model.TradeRequest) error {
	// Check crypto pair detail
	cryptoPairDetail, err := u.orderRepository.GetPairDetailByID(ctx, tradeReq.PairID)
	if err != nil {
		return err
	}

	// Lock taker and maker order id
	mutexMaker := u.redisLock.NewMutex(fmt.Sprintf("locking#maker#%v", tradeReq.MakerOrderID))
	if err := mutexMaker.Lock(); err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	mutexTaker := u.redisLock.NewMutex(fmt.Sprintf("locking#taker#%v", tradeReq.TakerOrderID))
	if err := mutexTaker.Lock(); err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	defer func() {
		if ok, err := mutexMaker.Unlock(); !ok || err != nil {
			log.Context(ctx).Error(err)
		}
		if ok, err := mutexTaker.Unlock(); !ok || err != nil {
			log.Context(ctx).Error(err)
		}
	}()

	// Get order detail from maker and taker
	takerOrder, err := u.orderRepository.GetOrder(ctx, tradeReq.TakerOrderID)
	if err != nil {
		return err
	}

	makerOrder, err := u.orderRepository.GetOrder(ctx, tradeReq.MakerOrderID)
	if err != nil {
		return err
	}

	takerOrder.FilledQuantity += tradeReq.Quantity
	makerOrder.FilledQuantity += tradeReq.Quantity

	// Partial filled
	takerOrder.Status = model.OrderStatusPartial
	makerOrder.Status = model.OrderStatusPartial

	// Update status to complete filled
	if takerOrder.FilledQuantity == takerOrder.Quantity {
		takerOrder.Status = model.OrderStatusComplete
	}
	if makerOrder.FilledQuantity == makerOrder.Quantity {
		makerOrder.Status = model.OrderStatusComplete
	}

	ctx, tx, err := gormpkg.InitTransactionToContext(ctx, u.writeDB)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	switch tradeReq.Side {
	case model.OrderSideBuy:
		// Update taker (buyer) primary pair wallet
		if err = u.orderRepository.UpdateUserWallet(ctx, takerOrder.UserID, cryptoPairDetail.PrimaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

		// Update maker (seller) secondary pair wallet
		if err = u.orderRepository.UpdateUserWallet(ctx, makerOrder.UserID, cryptoPairDetail.SecondaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

	case model.OrderSideSell:
		// Update taker (seller) secondary pair wallet
		if err = u.orderRepository.UpdateUserWallet(ctx, takerOrder.UserID, cryptoPairDetail.SecondaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

		// Update maker (buyer) primary pair wallet
		if err = u.orderRepository.UpdateUserWallet(ctx, makerOrder.UserID, cryptoPairDetail.PrimaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}
	}

	// Update order status transaction
	if _, err := u.orderRepository.SaveOrder(ctx, takerOrder); err != nil {
		return err
	}

	if _, err := u.orderRepository.SaveOrder(ctx, makerOrder); err != nil {
		return err
	}

	// Save to table match order
	matchOrder := model.MatchOrder{
		PairID:          tradeReq.PairID,
		TakerOrderID:    tradeReq.TakerOrderID,
		MakerOrderID:    tradeReq.MakerOrderID,
		Quantity:        tradeReq.Quantity,
		Price:           tradeReq.Price,
		TransactionTime: tradeReq.TradeTime,
	}

	if err = u.orderRepository.SaveMatchOrder(ctx, matchOrder); err != nil {
		return err
	}

	return tx.Commit().Error
}
