package main

import "fmt"

type ActionType string
const (
	ACTION_BID = "BID"
	ACTION_ASK = "ASK"
	ACTION_CANCEL = "CANCEL"
	ACTION_EXECUTE = "EXECUTION"
)

type Action struct {
	action_type ActionType
	symbol string
	order_id OrderID
	order_id_other OrderID
	price Price
	size Size
	trader TraderID
	trader_other TraderID
}

func newBidAction(order *Order) *Action {
    return &Action{
		action_type: ACTION_BID,
		symbol: order.symbol,
		order_id: order.order_id,
		price: order.price,
		size: order.size,
		trader: order.trader
	}
}

func newAskAction(order *Order) *Action {
    return &Action{
		action_type: ACTION_ASK,
		symbol: order.symbol,
		order_id: order.order_id,
		price: order.price,
		size: order.size,
		trader: order.trader
	}
}

func newCancelAction(order *Order) *Action {
    return &Action{
		action_type: ACTION_CANCEL,
		symbol: order.symbol,
		order_id: order.order_id,
		price: order.price,
		size: order.size,
		trader: order.trader
	}
}

func newExecuteAction(order *Order, entry *Order) *Action {
	// If statement to report execution based on side
    return &Action{
		action_type: ACTION_EXECUTE,
		symbol: order.symbol,
		order_id: order.order_id, // Change this around?
		order_id_other: entry.order_id, // Change this around?
		price: entry.price,
		size: order.size, // Change this around?
		trader: order.trader, // Change this around?
		trader_other: entry.trader // Change this around?
	}
}