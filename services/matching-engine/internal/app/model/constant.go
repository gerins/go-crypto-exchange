package model

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
