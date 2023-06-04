package model

import "encoding/json"

type Order struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	PairID          int     `json:"pair_id"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	Type            Type    `json:"type"`
	Side            Side    `json:"side"`
	Status          Status  `json:"status"`
	TransactionTime int64   `json:"transaction_time"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
