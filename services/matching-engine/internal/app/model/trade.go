package model

import "encoding/json"

const (
	BuyOrderCode  = "BUY"
	SellOrderCode = "SELL"
)

type Order struct {
	ID     string `json:"id"`
	Amount uint64 `json:"amount"`
	Price  uint64 `json:"price"`
	Type   string `json:"type"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
