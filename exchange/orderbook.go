package exchange

import (
	"sync"

	"github.com/gammazero/deque"
	"github.com/google/btree"
)

// PricePoint represents a price level and its associated orders Deque in the orderbook
type PricePoint struct {
	price  Price
	orders deque.Deque[OrderID]
	mutex  sync.Mutex
}

// Less is used by the btree package to compare PricePoints and allow nodes to be stored correctly
func (p *PricePoint) Less(than btree.Item) bool {
	return p.price < than.(*PricePoint).price
}

// OrderBook represents the collection of asks and bids, for a specific symbol on the exchange
type OrderBook struct {
	symbol   string
	asks     *btree.BTree
	bids     *btree.BTree
	exchange *Exchange
	mutex    sync.RWMutex
}

// init initialises the OrderBook with the given symbol and exchange and creates the btrees
func (ob *OrderBook) init(symbol string, exchange *Exchange) {
	// Lock the orderbook mutex to prevent concurrent access
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	ob.symbol = symbol
	ob.exchange = exchange

	ob.asks = btree.New(int(MAX_PRICE))
	ob.bids = btree.New(int(MAX_PRICE))
}

// limitHandle processes an incoming order in the following manner:
// 1. Immediately try to fill the incoming order
// 2. If the order is unfilled or partially filled, insert it into the orderbook
func (ob *OrderBook) limitHandle(incoming_order Order) {
	// Lock the orderbook mutex to prevent concurrent access
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

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

		// Lock the price point mutex to prevent concurrent access
		pp.mutex.Lock()
		defer pp.mutex.Unlock()

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

		// Lock the price point mutex to prevent concurrent access
		pp.mutex.Lock()
		defer pp.mutex.Unlock()

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
func (ob *OrderBook) fillOrder(order *Order, entries *deque.Deque[OrderID]) {
	// Lock the exchange mutex to prevent concurrent access
	ob.exchange.mutex.Lock()
	defer ob.exchange.mutex.Unlock()

	// Look up the order in the orderIDMap by the order_id
	if entry, ok := ob.exchange.orderIDMap[entries.Front()]; ok {
		// The existing book order is larger than the incoming order
		// Therefore, the incoming order is completely filled
		if entry.size >= order.size {
			// Report the trade to the exchange via the actions channel
			ob.exchange.actions <- newExecuteAction(order, &entry, order.size)

			// Reduce the existing book order size by the incoming order size and update the orderIDMap
			entry.size -= order.size
			ob.exchange.orderIDMap[entries.Front()] = entry

			// Reduce the incoming order size to zero to show that no further trades are possible
			order.size = 0
		} else {
			// The existing book order is smaller than the incoming order
			// Therefore, the incoming order is partially filled

			// Skip and remove cancelled orders (which have a size of zero from the cancel function)
			if entry.size == 0 {
				entries.PopFront()
				return
			}

			// Report the trade to the exchange via the actions channel
			ob.exchange.actions <- newExecuteAction(order, &entry, entry.size)

			// Reduce the incoming order size by the existing book order size
			order.size -= entry.size

			// Remove the existing book order from the orderbook and orderIDMap
			entries.PopFront()
			delete(ob.exchange.orderIDMap, entry.order_id)
		}
	} else {
		// The order_id is cannot be found in the order_id_map, so remove it from the orderbook
		entries.PopFront()
	}
}

// insertIntoBook inserts an incoming order into the appropriate btree of the orderbook
func (ob *OrderBook) insertIntoBook(order *Order) {

	// Select the appropriate btree based on the order side
	var tree *btree.BTree
	if order.side == Bid {
		tree = ob.bids
	} else {
		tree = ob.asks
	}

	// Create a new PricePoint with the order price
	pp := &PricePoint{price: order.price}

	// Check if the price point already exists in the orderbook
	if item := tree.Get(pp); item != nil {
		pp = item.(*PricePoint)
	}

	// Lock the price point mutex to prevent concurrent access
	pp.mutex.Lock()

	// Add the order to the price point's orders deque
	pp.orders.PushBack(order.order_id)

	// Unlock the price point mutex again
	pp.mutex.Unlock()

	// Insert the price point into the orderbook
	tree.ReplaceOrInsert(pp)

	// Lock the exchange mutex to prevent concurrent access
	ob.exchange.mutex.Lock()

	// Update the orderIDMap with the order details
	ob.exchange.orderIDMap[order.order_id] = *order

	// Unlock the exchange mutex again
	ob.exchange.mutex.Unlock()
}
