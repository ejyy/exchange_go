package exchange

import (
	"github.com/gammazero/deque"
	"github.com/google/btree"
)

// PricePoint represents a price level and its associated orders Deque in the orderbook
type PricePoint struct {
	price  Price
	orders deque.Deque[OrderID]
}

// Less is used by the btree package to compare PricePoints and allow nodes to be stored correctly
func (p *PricePoint) Less(than btree.Item) bool {
	return p.price < than.(*PricePoint).price
}

// OrderBook represents the orderbooks (asks and bids) for a specific symbol on the exchange
type OrderBook struct {
	symbol   string
	asks     *btree.BTree
	bids     *btree.BTree
	exchange *Exchange
}

// init initialises the OrderBook with the given symbol and exchange and creates the btrees
func (ob *OrderBook) init(symbol string, exchange *Exchange) {
	ob.symbol = symbol
	ob.exchange = exchange

	ob.asks = btree.New(int(MAX_PRICE))
	ob.bids = btree.New(int(MAX_PRICE))
}

// limitHandle processes an incoming order in the following manner:
// 1. Immediately try to fill the incoming order
// 2. If the order is unfilled or partially filled, insert it into the orderbook
func (ob *OrderBook) limitHandle(incoming_order Order) {
	order := incoming_order

	// Report the incoming order to the exchange via the actions channel
	ob.exchange.actions <- newOrderAction(&order)

	// Try to immediately fill the incoming order
	if order.side == Bid {
		ob.fillBidSide(&order)
	} else {
		ob.fillAskSide(&order)
	}

	// If unfilled (or partially filled), insert into the orderbook
	if order.size > 0 {
		ob.insertIntoBook(&order)
	}
}

// fillBidSide attempts to fill an incoming bid order by matching it with the lowest ask prices
func (ob *OrderBook) fillBidSide(order *Order) {
	// Find the minimum ask price that matches the incoming bid
	minAsk := ob.asks.Min()
	if minAsk == nil || order.price < minAsk.(*PricePoint).price {
		return // No matching asks
	}

	// Iterate through the existing book asks from lowest to highest price
	ob.asks.AscendGreaterOrEqual(minAsk, func(i btree.Item) bool {
		pp := i.(*PricePoint)

		// If the incoming bid price is less than the current ask price, stop iterating
		if order.price < pp.price || order.size == 0 {
			return false
		}

		// Fill the incoming bid order with the existing book asks
		for pp.orders.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, &pp.orders)
		}

		// If the price point is empty, remove it from the orderbook
		if pp.orders.Len() == 0 {
			ob.asks.Delete(pp)
		} else {
			// Otherwise, replace the price point in the orderbook
			ob.asks.ReplaceOrInsert(pp)
		}
		return true
	})
}

// fillAskSide attempts to fill an incoming ask order by matching it with the highest bid prices
func (ob *OrderBook) fillAskSide(order *Order) {
	// Find the maximum bid price that matches the incoming ask
	maxBid := ob.bids.Max()
	if maxBid == nil || order.price > maxBid.(*PricePoint).price {
		return // No matching bids
	}

	// Iterate through the existing book bids from highest to lowest price
	ob.bids.DescendLessOrEqual(maxBid, func(i btree.Item) bool {
		pp := i.(*PricePoint)

		// If the incoming ask price is greater than the current bid price, stop iterating
		if order.price > pp.price || order.size == 0 {
			return false
		}

		// Fill the incoming ask order with the existing book bids
		for pp.orders.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, &pp.orders)
		}

		// If the price point is empty, remove it from the orderbook
		if pp.orders.Len() == 0 {
			ob.bids.Delete(pp)
		} else {
			// Otherwise, replace the price point in the orderbook
			ob.bids.ReplaceOrInsert(pp)
		}
		return true
	})
}

// fillOrder fills an incoming order with the existing book orders
// TODO: ********** RESTART DOCUMENTATION FROM HERE **********
func (ob *OrderBook) fillOrder(order *Order, entries *deque.Deque[OrderID]) {
	if entry, ok := ob.exchange.order_id_map[entries.Front()]; ok {
		if entry.size >= order.size { // Incoming order completely filled
			ob.exchange.actions <- newExecuteAction(order, &entry, order.size)
			entry.size -= order.size
			ob.exchange.order_id_map[entries.Front()] = entry
			order.size = 0
		} else { // Incoming order partially filled
			// Skip cancelled orders
			if entry.size == 0 {
				entries.PopFront()
				return
			}
			ob.exchange.actions <- newExecuteAction(order, &entry, entry.size)
			order.size -= entry.size
			entries.PopFront()
			delete(ob.exchange.order_id_map, entry.order_id)
		}
	} else {
		entries.PopFront() // If order_id not found, potentially corrupted order so remove from orderbook
	}
}

func (ob *OrderBook) insertIntoBook(order *Order) {
	// Highlight: Insert into the appropriate btree
	var tree *btree.BTree
	if order.side == Bid {
		tree = ob.bids
	} else {
		tree = ob.asks
	}

	pp := &PricePoint{price: order.price}
	if item := tree.Get(pp); item != nil {
		pp = item.(*PricePoint)
	}
	pp.orders.PushBack(order.order_id)
	tree.ReplaceOrInsert(pp)

	ob.exchange.order_id_map[order.order_id] = *order
}
