package model

import "context"

type Usecase interface {
	ProcessOrder(ctx context.Context, orderReq RequestOrder) (Order, error)
}

type Repository interface {
	SaveOrder(ctx context.Context, order Order) (Order, error)
	GetPairDetail(ctx context.Context, code string) (Pair, error)
	GetUserWallet(ctx context.Context, userID, cryptoID int) (Wallet, error)
}
