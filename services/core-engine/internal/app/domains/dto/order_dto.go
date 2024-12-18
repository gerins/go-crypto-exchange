package dto

import "encoding/json"

type OrderRequest struct {
	PairCode string  `json:"pair_code"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Side     string  `json:"side"` // BUY / SELL
	Type     string  `json:"type"` // MARKET / LIMIT
}

type TradeRequest struct {
	PairID       int     `json:"pair_id"`
	TakerUserID  int     `json:"taker_user_id"`
	TakerOrderID int     `json:"taker_order_id"`
	MakerUserID  int     `json:"maker_user_id"`
	MakerOrderID int     `json:"maker_order_id"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	Side         string  `json:"side"`
	TradeTime    int64   `json:"trade_time"`
}

func (trade *TradeRequest) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

type BulkTradeRequest []TradeRequest

func (trade *BulkTradeRequest) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}
