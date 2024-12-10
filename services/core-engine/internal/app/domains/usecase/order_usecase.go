package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"core-engine/internal/app/domains/dto"
	"core-engine/internal/app/domains/model"
	serverError "core-engine/pkg/error"
	gormpkg "core-engine/pkg/gorm"
	"core-engine/pkg/jwt"
	"core-engine/pkg/kafka"
)

type orderUsecase struct {
	redisLock        *redsync.Redsync
	writeDB          *gorm.DB
	kafkaProducer    kafka.Producer
	validator        *validator.Validate
	orderRepository  model.OrderRepository
	userRepository   model.UserRepository
	walletRepository model.WalletRepository
}

// NewOrderUsecase returns new order usecase.
func NewOrderUsecase(
	writeDB *gorm.DB,
	redisLock *redsync.Redsync,
	kafkaProducer kafka.Producer,
	validator *validator.Validate,
	orderRepository model.OrderRepository,
	userRepository model.UserRepository,
	walletRepository model.WalletRepository,
) *orderUsecase {
	return &orderUsecase{
		writeDB:          writeDB,
		redisLock:        redisLock,
		kafkaProducer:    kafkaProducer,
		validator:        validator,
		orderRepository:  orderRepository,
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (u *orderUsecase) ProcessOrder(ctx context.Context, orderReq dto.OrderRequest) (model.Order, error) {
	defer log.Context(ctx).RecordDuration("ProcessOrder").Stop()

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
	cryptoPairDetail, err := u.walletRepository.GetPairDetail(ctx, orderReq.PairCode)
	if err != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	targetCryptoID := cryptoPairDetail.PrimaryCryptoID
	if model.Side(orderReq.Side) == model.OrderSideBuy {
		// When buying, check if user have enough secondary balance for buying primary crypto
		targetCryptoID = cryptoPairDetail.SecondaryCryptoID
	}

	// Lock all balance activity for this specific user
	timeRecord := log.Context(ctx).RecordDuration("obtaining lock")
	mutex := u.redisLock.NewMutex(fmt.Sprintf("locking#member#%v#%v", userDetail.ID, targetCryptoID))
	if err := mutex.Lock(); err != nil {
		log.Context(ctx).Error(err)
		return model.Order{}, err
	}

	timeRecord.Stop()

	defer func() { // Release the lock so other processes or threads can obtain a lock.
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Context(ctx).Error(err)
		}
	}()

	userWallet, err := u.walletRepository.GetUserWallet(ctx, userDetail.ID, targetCryptoID)
	if err != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	// Validate user balance
	if !userWallet.IsEnoughBalance(orderReq) {
		return model.Order{}, serverError.ErrInsufficientBalance(nil)
	}

	ctx, tx := gormpkg.InitTransactionToContext(ctx, u.writeDB)
	defer tx.WithContext(ctx).Rollback()

	// Deduct user wallet balance
	var errBalanceUpdate error
	switch model.Side(orderReq.Side) {
	case model.OrderSideSell:
		errBalanceUpdate = u.walletRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -orderReq.Quantity)

	case model.OrderSideBuy:
		totalAmount := orderReq.Price * orderReq.Quantity
		errBalanceUpdate = u.walletRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -totalAmount)
	}

	if errBalanceUpdate != nil {
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	newOrder := model.Order{
		UserID:          userDetail.ID,
		PairID:          cryptoPairDetail.ID,
		Quantity:        orderReq.Quantity,
		Price:           orderReq.Price,
		Type:            model.Type(orderReq.Type),
		Side:            model.Side(orderReq.Side),
		Status:          model.OrderStatusProgress,
		TransactionTime: time.Now().Unix(),
	}

	// Save to table orders
	order, err := u.orderRepository.SaveOrder(ctx, newOrder)
	if err != nil {
		return model.Order{}, err
	}

	// Publish to matching engine
	if err := u.kafkaProducer.Send(ctx, cryptoPairDetail.Code, cast.ToString(order.ID), order); err != nil {
		return model.Order{}, err
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Order{}, serverError.ErrGeneralDatabaseError(err)
	}

	return order, nil
}

func (u *orderUsecase) MatchOrder(ctx context.Context, tradeReq dto.TradeRequest) error {
	defer log.Context(ctx).RecordDuration("MatchOrder").Stop()

	// Check crypto pair detail
	cryptoPairDetail, err := u.walletRepository.GetPairDetailByID(ctx, tradeReq.PairID)
	if err != nil {
		return err
	}

	// Compose locking key
	lockCombination := []int{tradeReq.MakerUserID, tradeReq.TakerUserID}

	// Deterministic order
	// This approach prevents deadlocks, which typically occur when transactions acquire locks on
	// the same resources but in a different order.
	sort.Ints(lockCombination)

	timeRecord := log.Context(ctx).RecordDuration("obtaining lock")
	lock := u.redisLock.NewMutex(fmt.Sprintf("locking#trade#%v", lockCombination))
	if err := lock.Lock(); err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	timeRecord.Stop()

	defer func() {
		if ok, err := lock.Unlock(); !ok || err != nil {
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

	ctx, tx := gormpkg.InitTransactionToContext(ctx, u.writeDB)
	defer tx.WithContext(ctx).Rollback()

	switch model.Side(tradeReq.Side) {
	case model.OrderSideBuy:
		// Update taker (buyer) primary pair wallet
		if err = u.walletRepository.UpdateUserWallet(ctx, takerOrder.UserID, cryptoPairDetail.PrimaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

		// Update maker (seller) secondary pair wallet
		if err = u.walletRepository.UpdateUserWallet(ctx, makerOrder.UserID, cryptoPairDetail.SecondaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

	case model.OrderSideSell:
		// Update taker (seller) secondary pair wallet
		if err = u.walletRepository.UpdateUserWallet(ctx, takerOrder.UserID, cryptoPairDetail.SecondaryCryptoID, tradeReq.Quantity); err != nil {
			return err
		}

		// Update maker (buyer) primary pair wallet
		if err = u.walletRepository.UpdateUserWallet(ctx, makerOrder.UserID, cryptoPairDetail.PrimaryCryptoID, tradeReq.Quantity); err != nil {
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

	return tx.WithContext(ctx).Commit().Error
}
