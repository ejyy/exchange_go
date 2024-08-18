package exchange

import (
	"github.com/gammazero/deque"
	"github.com/google/btree"
)

// Highlight: New struct to represent a price point in the btree
type PricePoint struct {
	price  Price
	orders deque.Deque[OrderID]
}

// Highlight: Implement Less interface for btree.Item
func (p *PricePoint) Less(than btree.Item) bool {
	return p.price < than.(*PricePoint).price
}

type OrderBook struct {
	symbol string
	// Highlight: Replace price_points array with two btrees
	asks     *btree.BTree
	bids     *btree.BTree
	exchange *Exchange
}

func (ob *OrderBook) init(symbol string, exchange *Exchange) {
	ob.symbol = symbol
	ob.exchange = exchange

	// Highlight: Initialize btrees
	ob.asks = btree.New(100_000)
	ob.bids = btree.New(100_000)
}

func (ob *OrderBook) limitHandle(incoming_order Order) {
	order := incoming_order

	ob.exchange.actions <- newOrderAction(&order)

	// Try to immediately fill the incoming order
	if order.side == Bid {
		ob.fillBidSide(&order)
	} else {
		ob.fillAskSide(&order)
	}

	// If unfilled (or partially filled), insert into orderbook
	if order.size > 0 {
		ob.insertIntoBook(&order)
	}
}

func (ob *OrderBook) fillBidSide(order *Order) {
	// Find the minimum ask price that matches the incoming bid
	minAsk := ob.asks.Min()
	if minAsk == nil || order.price < minAsk.(*PricePoint).price {
		return // No matching asks
	}

	ob.asks.AscendGreaterOrEqual(minAsk, func(i btree.Item) bool {
		pp := i.(*PricePoint)
		if order.price < pp.price || order.size == 0 {
			return false
		}
		for pp.orders.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, &pp.orders)
		}
		if pp.orders.Len() == 0 {
			ob.asks.Delete(pp)
		} else {
			ob.asks.ReplaceOrInsert(pp)
		}
		return true
	})
}

func (ob *OrderBook) fillAskSide(order *Order) {
	// Find the maximum bid price that matches the incoming ask
	maxBid := ob.bids.Max()
	if maxBid == nil || order.price > maxBid.(*PricePoint).price {
		return // No matching bids
	}

	ob.bids.DescendLessOrEqual(maxBid, func(i btree.Item) bool {
		pp := i.(*PricePoint)
		if order.price > pp.price || order.size == 0 {
			return false
		}
		for pp.orders.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, &pp.orders)
		}
		if pp.orders.Len() == 0 {
			ob.bids.Delete(pp)
		} else {
			ob.bids.ReplaceOrInsert(pp)
		}
		return true
	})
}

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

// Highlight: Remove updateBidMaxAskMin function as it's no longer needed
