package model

type RequestOrder struct {
	PairCode string  `json:"pair_code"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Side     Side    `json:"side"` // BUY / SELL
	Type     Type    `json:"type"` // MARKET / LIMIT
}
