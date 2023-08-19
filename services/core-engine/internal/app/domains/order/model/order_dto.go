package model

import "encoding/json"

type OrderRequest struct {
	PairCode string  `json:"pair_code"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Side     Side    `json:"side"` // BUY / SELL
	Type     Type    `json:"type"` // MARKET / LIMIT
}

type TradeRequest struct {
	PairID       int     `json:"pair_id"`
	TakerOrderID int     `json:"taker_order_id"`
	MakerOrderID int     `json:"maker_order_id"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	Side         Side    `json:"side"`
	TradeTime    int64   `json:"trade_time"`
}

type BulkTradeRequest []TradeRequest

func (trade *BulkTradeRequest) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}
