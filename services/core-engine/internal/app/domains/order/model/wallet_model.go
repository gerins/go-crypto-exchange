package model

import "time"

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

func (userWallet Wallet) IsEnoughBalance(orderReq OrderRequest) bool {
	switch orderReq.Side {
	case OrderSideSell:
		return userWallet.Quantity >= orderReq.Quantity

	case OrderSideBuy:
		totalBuyAmount := orderReq.Price * orderReq.Quantity
		return userWallet.Quantity >= totalBuyAmount
	}

	return false
}
