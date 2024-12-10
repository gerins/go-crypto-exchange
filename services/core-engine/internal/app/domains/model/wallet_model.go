package model

import (
	"context"
	"time"

	"core-engine/internal/app/domains/dto"
)

type Wallet struct {
	ID        int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID    int        `json:"user_id" gorm:"column:user_id;type:int"`
	CryptoID  int        `json:"crypto_id" gorm:"column:crypto_id;type:int"`
	Quantity  float64    `json:"quantity" gorm:"column:quantity;type:text"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
}

func (Wallet) TableName() string {
	return "wallet"
}

func (userWallet Wallet) IsEnoughBalance(orderReq dto.OrderRequest) bool {
	switch Side(orderReq.Side) {
	case OrderSideSell:
		return userWallet.Quantity >= orderReq.Quantity

	case OrderSideBuy:
		totalBuyAmount := orderReq.Price * orderReq.Quantity
		return userWallet.Quantity >= totalBuyAmount
	}

	return false
}

type WalletRepository interface {
	// Crypto Pair
	GetPairDetail(ctx context.Context, code string) (Pair, error)
	GetPairDetailByID(ctx context.Context, id int) (Pair, error)

	// Wallet
	Save(ctx context.Context, wallet Wallet) error
	GetUserWallet(ctx context.Context, userID, cryptoID int) (Wallet, error)
	UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error
}
