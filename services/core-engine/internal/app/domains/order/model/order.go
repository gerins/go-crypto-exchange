package model

import "context"

type Usecase interface {
	ProcessOrder(ctx context.Context, orderReq OrderRequest) (Order, error)
	MatchOrder(ctx context.Context, tradeReq TradeRequest) error
}

type Repository interface {
	// Crypto Pair
	GetPairDetail(ctx context.Context, code string) (Pair, error)
	GetPairDetailByID(ctx context.Context, id int) (Pair, error)

	// User Order
	SaveOrder(ctx context.Context, order Order) (Order, error)
	GetOrder(ctx context.Context, id int) (Order, error)

	// Matching Order
	SaveMatchOrder(ctx context.Context, matchOrder MatchOrder) error

	// Wallet
	GetUserWallet(ctx context.Context, userID, cryptoID int) (Wallet, error)
	UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error
}
