package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"

	"matching-engine/internal/app/model"
	"matching-engine/pkg/kafka"
)

// orderBook is used for processing data orderBook
type OrderBook struct {
	matchOrderTopic string
	cache           *redis.Client
	kafkaProducer   kafka.Producer
	validator       *validator.Validate
	BuyOrders       []model.Order
	SellOrders      []model.Order
}

// NewOrderBook returns new order book usecase.
func NewOrderBook(
	matchOrderTopic string,
	cache *redis.Client,
	kafkaProducer kafka.Producer,
	validator *validator.Validate,
) *OrderBook {
	return &OrderBook{
		matchOrderTopic: matchOrderTopic,
		cache:           cache,
		kafkaProducer:   kafkaProducer,
		validator:       validator,
		BuyOrders:       []model.Order{},
		SellOrders:      []model.Order{},
	}
}

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Execute(ctx context.Context, order model.Order) error {
	var trades []model.Trade

	switch order.Side {
	case model.OrderSideBuy:
		trades = book.processLimitBuy(order)

	case model.OrderSideSell:
		trades = book.processLimitSell(order)
	}

	if len(trades) == 0 {
		return nil
	}

	log.Context(ctx).RespBody = trades

	// Publish to Kafka
	if err := book.kafkaProducer.Send(ctx, book.matchOrderTopic, cast.ToString(order.ID), trades); err != nil {
		return err
	}

	return nil
}

// Process a limit buy order
func (book *OrderBook) processLimitBuy(reqOrder model.Order) []model.Trade {
	trades := make([]model.Trade, 0, 1)
	n := len(book.SellOrders)

	// check if we have at least one matching order
	if n != 0 && book.SellOrders[n-1].Price <= reqOrder.Price {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price > reqOrder.Price {
				break
			}
			// fill the entire order
			if sellOrder.Quantity >= reqOrder.Quantity {
				tradeTime := time.Now().Unix()
				trades = append(trades, model.Trade{PairID: reqOrder.PairID, TakerOrderID: reqOrder.ID, MakerOrderID: sellOrder.ID, Quantity: reqOrder.Quantity, Price: sellOrder.Price, TradeTime: tradeTime})
				sellOrder.Quantity -= reqOrder.Quantity
				if sellOrder.Quantity == 0 {
					book.removeSellOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if sellOrder.Quantity < reqOrder.Quantity {
				tradeTime := time.Now().Unix()
				trades = append(trades, model.Trade{PairID: reqOrder.PairID, TakerOrderID: reqOrder.ID, MakerOrderID: sellOrder.ID, Quantity: sellOrder.Quantity, Price: sellOrder.Price, TradeTime: tradeTime})
				reqOrder.Quantity -= sellOrder.Quantity
				book.removeSellOrder(i)
				continue
			}
		}
	}

	// finally add the remaining order to the list
	book.addBuyOrder(reqOrder)
	return trades
}

// Process a limit sell order
func (book *OrderBook) processLimitSell(reqOrder model.Order) []model.Trade {
	trades := make([]model.Trade, 0, 1)
	n := len(book.BuyOrders)
	// check if we have at least one matching order
	if n != 0 && book.BuyOrders[n-1].Price >= reqOrder.Price {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price < reqOrder.Price {
				break
			}
			// fill the entire order
			if buyOrder.Quantity >= reqOrder.Quantity {
				tradeTime := time.Now().Unix()
				trades = append(trades, model.Trade{PairID: reqOrder.PairID, TakerOrderID: reqOrder.ID, MakerOrderID: buyOrder.ID, Quantity: reqOrder.Quantity, Price: buyOrder.Price, TradeTime: tradeTime})
				buyOrder.Quantity -= reqOrder.Quantity
				if buyOrder.Quantity == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			if buyOrder.Quantity < reqOrder.Quantity {
				tradeTime := time.Now().Unix()
				trades = append(trades, model.Trade{PairID: reqOrder.PairID, TakerOrderID: reqOrder.ID, MakerOrderID: buyOrder.ID, Quantity: buyOrder.Quantity, Price: buyOrder.Price, TradeTime: tradeTime})
				reqOrder.Quantity -= buyOrder.Quantity
				book.removeBuyOrder(i)
				continue
			}
		}
	}

	// finally add the remaining order to the list
	book.addSellOrder(reqOrder)
	return trades
}

// Add a buy order to the order book
// Buy orders is arrange in cheapest...most expensive
func (book *OrderBook) addBuyOrder(order model.Order) {
	// Find the index where the value should be inserted
	index := sort.Search(len(book.BuyOrders), func(i int) bool {
		return book.BuyOrders[i].Price >= order.Price
	})

	book.BuyOrders = append(book.BuyOrders[:index], append([]model.Order{order}, book.BuyOrders[index:]...)...)
}

// Add a sell order to the order book
// Sell orders is arrange in most expensive...cheapest
func (book *OrderBook) addSellOrder(order model.Order) {
	// Find the index where the value should be inserted
	index := sort.Search(len(book.SellOrders), func(i int) bool {
		return book.SellOrders[i].Price <= order.Price
	})

	book.SellOrders = append(book.SellOrders[:index], append([]model.Order{order}, book.SellOrders[index:]...)...)
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
