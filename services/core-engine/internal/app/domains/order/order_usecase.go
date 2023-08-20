package order

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"

	"core-engine/internal/app/domains/order/model"
	"core-engine/internal/app/domains/user"
	"core-engine/pkg/jwt"
	"core-engine/pkg/kafka"
)

type usecase struct {
	kafkaProducer   kafka.Producer
	validator       *validator.Validate
	orderRepository model.Repository
	userRepository  user.Repository
}

// NewUsecase returns new order usecase.
func NewUsecase(
	kafkaProducer kafka.Producer,
	validator *validator.Validate,
	orderRepository model.Repository,
	userRepository user.Repository,
) model.Usecase {
	return &usecase{
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
		return model.Order{}, err
	}

	// Check account status
	if !userDetail.Status {
		return model.Order{}, user.ErrUserBlocked // User already deactivated
	}

	// Check crypto pair detail
	cryptoPairDetail, err := u.orderRepository.GetPairDetail(ctx, orderReq.PairCode)
	if err != nil {
		return model.Order{}, err
	}

	targetCryptoID := cryptoPairDetail.PrimaryCryptoID
	if orderReq.Side == model.OrderSideBuy {
		// When buying, check if user have enough secondary balance for buying primary crypto
		targetCryptoID = cryptoPairDetail.SecondaryCryptoID
	}

	userWallet, err := u.orderRepository.GetUserWallet(ctx, userDetail.ID, targetCryptoID)
	if err != nil {
		return model.Order{}, err
	}

	// Validate user balance
	if !userWallet.IsEnoughBalance(orderReq) {
		return model.Order{}, model.ErrInsufficientBalance
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

	// Deduct user wallet balance
	switch orderReq.Side {
	case model.OrderSideSell:
		u.orderRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -orderReq.Quantity)
	case model.OrderSideBuy:
		totalAmount := orderReq.Price * orderReq.Quantity
		u.orderRepository.UpdateUserWallet(ctx, userDetail.ID, userWallet.CryptoID, -totalAmount)
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

	// TODO : Add database transaction
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

	if u.orderRepository.SaveMatchOrder(ctx, matchOrder); err != nil {
		return err
	}

	return nil
}
