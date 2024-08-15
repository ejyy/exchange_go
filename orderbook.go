package main

import (
	// "fmt"
	"github.com/gammazero/deque"
)

type OrderBook struct {
	symbol           string
	current_order_id OrderID
	ask_min          Price
	bid_max          Price
	order_id_map     map[OrderID]Order // TODO: Wasteful to store entire 'Order' struct, only need trader + size (https://go.dev/play/p/4KPix5OEXJC)
	price_points     [MAX_PRICE + 1]deque.Deque[OrderID]
	// Consider actions channel here to handle message passing
}

func (ob *OrderBook) init(symbol string) {
	ob.symbol = symbol
	ob.current_order_id = 0

	ob.ask_min = MAX_PRICE + 1
	ob.bid_max = MIN_PRICE - 1

	ob.order_id_map = make(map[OrderID]Order, EST_ORDERS)
	for i := range MAX_PRICE + 1 {
		ob.price_points[i] = *deque.New[OrderID]()
	}

	// fmt.Println("Orderbook created for:", ob.symbol)
}

func (ob *OrderBook) limitHandle(incoming_order Order) OrderID {
	order := incoming_order

	ob.current_order_id += 1
	order.order_id = ob.current_order_id

	// fmt.Printf("ORDER... ID: %v, Symbol: %v, Side: %v, Price: %v, Size: %v, Trader: %v\n",
	// 	order.order_id, order.symbol, order.side, order.price, order.size, order.trader)

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

	return order.order_id
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

// func createExecution(order *Order, entry *Order, fill_size Size) Execution {
//     if order.side == Bid {
//         return Execution{
//             order_id_bid: order.order_id,
//             order_id_ask: entry.order_id,
//             price:        entry.price,
//             size:         fill_size,
//             trader_bid:   order.trader,
//             trader_ask:   entry.trader,
//         }
//     } else {
//         return Execution{
//             order_id_bid: entry.order_id,
//             order_id_ask: order.order_id,
//             price:        entry.price,
//             size:         fill_size,
//             trader_bid:   entry.trader,
//             trader_ask:   order.trader,
//         }
//     }
// }

// TODO: Tidy up execution reporting to minimise code repetition (suggest using a 'fillsize' call to a reporting function)
// Execution occurs at entry.price for 'price improvement'

func (ob *OrderBook) fillOrder(order *Order, entries *deque.Deque[OrderID]) {
	if entry, ok := ob.order_id_map[entries.Front()]; ok {
		if entry.size >= order.size { // Incoming order completely filled

			// if order.side == Bid {
			// 	fmt.Println(&Execution{symbol: order.symbol, order_id_bid: order.order_id, order_id_ask: entry.order_id,
			// 		price: entry.price, size: order.size, trader_bid: order.trader, trader_ask: entry.trader})
			// } else {
			// 	fmt.Println(&Execution{symbol: order.symbol, order_id_bid: entry.order_id, order_id_ask: order.order_id,
			// 		price: entry.price, size: order.size, trader_bid: entry.trader, trader_ask: order.trader})
			// }

			entry.size -= order.size
			ob.order_id_map[entries.Front()] = entry

			order.size = 0
		} else { // Incoming order partially filled

			// Skip cancelled orders
			if entry.size == 0 {
				entries.PopFront()
				return
			}

			// if order.side == Bid {
			// 	fmt.Println(&Execution{symbol: order.symbol, order_id_bid: order.order_id, order_id_ask: entry.order_id,
			// 		price: entry.price, size: entry.size, trader_bid: order.trader, trader_ask: entry.trader})
			// } else {
			// 	fmt.Println(&Execution{symbol: order.symbol, order_id_bid: entry.order_id, order_id_ask: order.order_id,
			// 		price: entry.price, size: entry.size, trader_bid: entry.trader, trader_ask: order.trader})
			// }

			order.size -= entry.size
			entries.PopFront()
			delete(ob.order_id_map, entry.order_id)
		}
	} else {
		entries.PopFront() // If order_id not found, potentially corrupted order so remove from orderbook
	}
}

func (ob *OrderBook) insertIntoBook(order *Order) {
	ob.price_points[order.price].PushBack(ob.current_order_id)
	ob.order_id_map[order.order_id] = *order

	ob.updateBidMaxAskMin(order)
}

func (ob *OrderBook) updateBidMaxAskMin(order *Order) {
	if order.side == Bid && order.price > ob.bid_max {
		ob.bid_max = order.price
	} else if order.side == Ask && order.price < ob.ask_min {
		ob.ask_min = order.price
	}
}
