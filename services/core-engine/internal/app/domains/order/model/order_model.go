package model

import (
	"errors"
	"time"
)

type (
	Side   string
	Type   string
	Status string
)

const (
	OrderSideBuy  Side = "BUY"
	OrderSideSell Side = "SELL"
)

const (
	OrderTypeMarket     Type = "MARKET"
	OrderTypeLimit      Type = "LIMIT"
	OrderTypeStopLoss   Type = "STOP_LOSS"
	OrderTypeTakeProfit Type = "TAKE_PROFIT"
)

const (
	OrderStatusComplete Status = "COMPLETE"
	OrderStatusFailed   Status = "FAILED"
	OrderStatusProgress Status = "PROGRESS"
	OrderStatusPartial  Status = "PARTIAL"
)

var (
	ErrInsufficientBalance = errors.New("Insufficient balance")
)

type Order struct {
	ID              int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID          int        `json:"user_id" gorm:"column:user_id;type:int"`
	PairID          int        `json:"pair_id" gorm:"column:pair_id;type:int"`
	Quantity        float64    `json:"quantity" gorm:"column:quantity;type:double"`
	FilledQuantity  float64    `json:"filled_quantity" gorm:"column:filled_quantity;type:double"`
	Price           float64    `json:"price" gorm:"column:price;type:double"`
	Type            Type       `json:"type" gorm:"column:type;type:text"`
	Side            Side       `json:"side" gorm:"column:side;type:text"`
	Status          Status     `json:"status" gorm:"column:status;type:text"`
	TransactionTime int64      `json:"transaction_time" gorm:"column:transaction_time;type:bigint"` // Transaction time
	CreatedAt       time.Time  `json:"-" gorm:"column:created_at;type:datetime"`
	UpdatedAt       time.Time  `json:"-" gorm:"column:updated_at;type:datetime"`
	DeletedAt       *time.Time `json:"-" gorm:"column:deleted_at;type:datetime"`
}

func (Order) TableName() string {
	return "orders"
}
