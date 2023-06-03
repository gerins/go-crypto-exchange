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

func (u *usecase) ProcessOrder(ctx context.Context, orderReq model.RequestOrder) (model.Order, error) {
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
