package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"

	"matching-engine/internal/app/model"
	"matching-engine/pkg/kafka"
)

// orderBook is used for processing data orderBook
type OrderBook struct {
	cache         *redis.Client
	kafkaProducer kafka.Producer
	validator     *validator.Validate
	BuyOrders     []model.Order
	SellOrders    []model.Order
}

// NewOrderBook returns new order book usecase.
func NewOrderBook(validator *validator.Validate, kafkaProducer kafka.Producer, cache *redis.Client) *OrderBook {
	return &OrderBook{
		cache:         cache,
		kafkaProducer: kafkaProducer,
		validator:     validator,
		BuyOrders:     []model.Order{},
		SellOrders:    []model.Order{},
	}
}

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Execute(ctx context.Context, order model.Order) []model.Trade {
	if order.Type == model.BuyOrderCode {
		return book.processLimitBuy(order)
	}

	return book.processLimitSell(order)
}

// Process a limit buy order
func (book *OrderBook) processLimitBuy(order model.Order) []model.Trade {
	trades := make([]model.Trade, 0, 1)
	n := len(book.SellOrders)
	// check if we have at least one matching order
	if n != 0 || book.SellOrders[n-1].Price <= order.Price {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price > order.Price {
				break
			}
			// fill the entire order
			if sellOrder.Amount >= order.Amount {
				trades = append(trades, model.Trade{TakerOrderID: order.ID, MakerOrderID: sellOrder.ID, Amount: order.Amount, Price: sellOrder.Price})
				sellOrder.Amount -= order.Amount
				if sellOrder.Amount == 0 {
					book.removeSellOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if sellOrder.Amount < order.Amount {
				trades = append(trades, model.Trade{TakerOrderID: order.ID, MakerOrderID: sellOrder.ID, Amount: sellOrder.Amount, Price: sellOrder.Price})
				order.Amount -= sellOrder.Amount
				book.removeSellOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addBuyOrder(order)
	return trades
}

// Process a limit sell order
func (book *OrderBook) processLimitSell(order model.Order) []model.Trade {
	trades := make([]model.Trade, 0, 1)
	n := len(book.BuyOrders)
	// check if we have at least one matching order
	if n != 0 || book.BuyOrders[n-1].Price >= order.Price {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price < order.Price {
				break
			}
			// fill the entire order
			if buyOrder.Amount >= order.Amount {
				trades = append(trades, model.Trade{TakerOrderID: order.ID, MakerOrderID: buyOrder.ID, Amount: order.Amount, Price: buyOrder.Price})
				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if buyOrder.Amount < order.Amount {
				trades = append(trades, model.Trade{TakerOrderID: order.ID, MakerOrderID: buyOrder.ID, Amount: buyOrder.Amount, Price: buyOrder.Price})
				order.Amount -= buyOrder.Amount
				book.removeBuyOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addSellOrder(order)
	return trades
}

// Add a buy order to the order book
func (book *OrderBook) addBuyOrder(order model.Order) {
	n := len(book.BuyOrders)
	var i int
	for i := n - 1; i >= 0; i-- {
		buyOrder := book.BuyOrders[i]
		if buyOrder.Price < order.Price {
			break
		}
	}
	if i == n-1 {
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		copy(book.BuyOrders[i+1:], book.BuyOrders[i:])
		book.BuyOrders[i] = order
	}
}

// Add a sell order to the order book
func (book *OrderBook) addSellOrder(order model.Order) {
	n := len(book.SellOrders)
	var i int
	for i := n - 1; i >= 0; i-- {
		sellOrder := book.SellOrders[i]
		if sellOrder.Price > order.Price {
			break
		}
	}
	if i == n-1 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		copy(book.SellOrders[i+1:], book.SellOrders[i:])
		book.SellOrders[i] = order
	}
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
