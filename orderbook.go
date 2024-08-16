package main

import (
	"github.com/gammazero/deque"
)

type OrderBook struct {
	symbol       string
	ask_min      Price
	bid_max      Price
	price_points [MAX_PRICE + 1]deque.Deque[OrderID]
	exchange     *Exchange
}

func (ob *OrderBook) init(symbol string, exchange *Exchange) {
	ob.symbol = symbol
	ob.exchange = exchange

	ob.ask_min = MAX_PRICE + 1
	ob.bid_max = MIN_PRICE - 1

	for i := range MAX_PRICE + 1 {
		ob.price_points[i] = *deque.New[OrderID]()
	}
}

func (ob *OrderBook) limitHandle(incoming_order Order) {
	order := incoming_order

	// Try to immediately fill the incoming order
	if order.side == Bid {
		ob.exchange.actions <- newOrderAction(&order)
		ob.fillBidSide(&order)
	} else {
		ob.exchange.actions <- newOrderAction(&order)
		ob.fillAskSide(&order)
	}

	// If unfilled (or partially filled), insert into orderbook
	if order.size > 0 {
		ob.insertIntoBook(&order)
	}
}

func (ob *OrderBook) fillBidSide(order *Order) {
	for order.price >= ob.ask_min && order.size > 0 {
		entries := &ob.price_points[ob.ask_min]
		for entries.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, entries)
		}
		if order.size > 0 {
			ob.ask_min++
		}
	}
}

func (ob *OrderBook) fillAskSide(order *Order) {
	for order.price <= ob.bid_max && order.size > 0 {
		entries := &ob.price_points[ob.bid_max]
		for entries.Len() > 0 && order.size > 0 {
			ob.fillOrder(order, entries)
		}
		if order.size > 0 {
			ob.bid_max--
		}
	}
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
	ob.price_points[order.price].PushBack(order.order_id)
	ob.exchange.order_id_map[order.order_id] = *order

	ob.updateBidMaxAskMin(order)
}

func (ob *OrderBook) updateBidMaxAskMin(order *Order) {
	if order.side == Bid && order.price > ob.bid_max {
		ob.bid_max = order.price
	} else if order.side == Ask && order.price < ob.ask_min {
		ob.ask_min = order.price
	}
}
