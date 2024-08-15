// Go exchange

package main

import (
	"fmt"
)

type Exchange struct {
	name           string
	orderbooks_map map[string]*OrderBook
	//actions chan<- *Action
}

// Consider: Could move current_order_id and order_id_map out of OrderBook and into Exchange to prevent having to re-make a map for each symbol

// TODO: Can pre-warm the Exchange by initialising with a list of symbols

func (ex *Exchange) Init(name string) {
	ex.name = name
	ex.orderbooks_map = make(map[string]*OrderBook, EST_SYMBOLS)
	//ex.actions = actions (need to make Init accept actions chan<- *Action)
	// fmt.Println("Exchange started:", ex.name, "- Ready to accept orders")
}

func (ex *Exchange) getOrCreateOrderBook(symbol string) *OrderBook {
	order_book, exists := ex.orderbooks_map[symbol]
	if !exists {
		order_book = new(OrderBook)
		order_book.init(symbol)
		ex.orderbooks_map[symbol] = order_book
	}
	return order_book
}

func (ex *Exchange) Limit(incoming_order Order) {
	// TODO: Edit this to do order validation

	ob := ex.getOrCreateOrderBook(incoming_order.symbol)
	ob.limitHandle(incoming_order)
}

func (ex *Exchange) Cancel(symbol string, order_id OrderID) OrderID {
	ob := ex.getOrCreateOrderBook(symbol)
	if cancel_order, ok := ob.order_id_map[order_id]; ok {
		if cancel_order.size == 0 {
			fmt.Println("CANCEL rejected as ID:", order_id, "already cancelled (or fully filled)")
		} else {
			(&cancel_order).size = 0
			ob.order_id_map[order_id] = cancel_order
			//ex.actions <- NewCancelAction(cancel_order)
			fmt.Println("CANCEL... ID:", order_id, "Symbol:", symbol)
		}
	} else {
		fmt.Println("CANCEL rejected as ID:", order_id, "not found")
	}
	return order_id
}
