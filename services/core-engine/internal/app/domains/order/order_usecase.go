package order

import (
	"context"

	"github.com/go-playground/validator/v10"

	"core-engine/internal/app/domains/order/model"
	"core-engine/internal/app/domains/user"
)

type usecase struct {
	validator       *validator.Validate
	orderRepository model.Repository
	userRepository  user.Repository
}

// NewUsecase returns new order usecase.
func NewUsecase(validator *validator.Validate, orderRepository model.Repository, userRepository user.Repository) model.Usecase {
	return &usecase{
		validator:       validator,
		orderRepository: orderRepository,
		userRepository:  userRepository,
	}
}

func (u *usecase) ProcessOrder(ctx context.Context, orderReq model.RequestOrder) (model.Order, error) {
	// Check user detail
	userDetail, err := u.userRepository.FindUserByEmail(ctx, "")
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
		UserID:   userDetail.ID,
		PairID:   cryptoPairDetail.ID,
		Quantity: orderReq.Quantity,
		Price:    orderReq.Price,
		Type:     orderReq.Type,
		Side:     orderReq.Side,
		Status:   model.OrderStatusProgress,
	}

	order, err := u.orderRepository.SaveOrder(ctx, newOrder)
	if err != nil {
		return model.Order{}, err
	}

	// Publish to matching engine

	return order, nil
}
