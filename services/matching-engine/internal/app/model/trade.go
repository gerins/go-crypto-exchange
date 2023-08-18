package model

import "encoding/json"

type Trade struct {
	PairID       int     `json:"pair_id"`
	TakerOrderID int     `json:"taker_order_id"`
	MakerOrderID int     `json:"maker_order_id"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	Side         Side    `json:"side"`
	TradeTime    int64   `json:"trade_time"`
}

func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}
