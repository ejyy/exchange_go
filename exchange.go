package main

import (
	"fmt"
)

type Exchange struct {
	name             string
	orderbooks_map   map[string]*OrderBook
	current_order_id OrderID
	order_id_map     map[OrderID]Order // TODO: Wasteful to store entire 'Order' struct, only need trader + size (https://go.dev/play/p/4KPix5OEXJC)
	actions          chan *Action
}

// TODO: Can pre-warm the Exchange by initialising with a list of symbols

func (ex *Exchange) Init(name string, actions chan *Action) {
	ex.name = name
	ex.current_order_id = 0

	ex.orderbooks_map = make(map[string]*OrderBook, EST_SYMBOLS)
	ex.order_id_map = make(map[OrderID]Order, EST_ORDERS)

	ex.actions = actions

	fmt.Println("Exchange started:", ex.name, "- Ready to accept orders")
}

func (ex *Exchange) getNextOrderID() OrderID {
	ex.current_order_id += 1
	return ex.current_order_id
}

func (ex *Exchange) getOrCreateOrderBook(symbol string) *OrderBook {
	order_book, exists := ex.orderbooks_map[symbol]
	if !exists {
		order_book = new(OrderBook)
		order_book.init(symbol, ex)
		ex.orderbooks_map[symbol] = order_book
	}
	return order_book
}

func (ex *Exchange) Limit(incoming_order Order) {
	// TODO: Add in order validation

	ob := ex.getOrCreateOrderBook(incoming_order.symbol)
	incoming_order.order_id = ex.getNextOrderID()
	ob.limitHandle(incoming_order)
}

func (ex *Exchange) Cancel(order_id OrderID) {
	if cancel_order, ok := ex.order_id_map[order_id]; ok {
		if cancel_order.size == 0 {
			ex.actions <- newCancelRejectAction()
		} else {
			cancel_order.size = 0
			ex.order_id_map[order_id] = cancel_order
			ex.actions <- newCancelAction(&cancel_order)
		}
	} else {
		ex.actions <- newCancelRejectAction()
	}
}
