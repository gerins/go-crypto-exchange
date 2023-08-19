package model

import "time"

type MatchOrder struct {
	ID              int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	PairID          int        `json:"pair_id" gorm:"column:pair_id;type:int"`
	TakerOrderID    int        `json:"taker_order_id" gorm:"column:taker_order_id;type:int"`
	MakerOrderID    int        `json:"maker_order_id" gorm:"column:maker_order_id;type:int"`
	Quantity        float64    `json:"quantity" gorm:"column:quantity;type:double"`
	Price           float64    `json:"price" gorm:"column:price;type:double"`
	TransactionTime int64      `json:"transaction_time" gorm:"column:transaction_time;type:bigint"` // Transaction time
	CreatedAt       time.Time  `json:"-" gorm:"column:created_at;type:datetime"`
	UpdatedAt       time.Time  `json:"-" gorm:"column:updated_at;type:datetime"`
	DeletedAt       *time.Time `json:"-" gorm:"column:deleted_at;type:datetime"`
}

func (MatchOrder) TableName() string {
	return "match_orders"
}
