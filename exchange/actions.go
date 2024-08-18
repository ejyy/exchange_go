package exchange

import "fmt"

// ActionType represents the type of action event passed by the exchange
type ActionType string

// Define the action types used in the exchange for various order event states
const (
	ACTION_BID           = "ORDER_BID"
	ACTION_ASK           = "ORDER_ASK"
	ACTION_ORDER_REJECT  = "ORDER_REJECT"
	ACTION_CANCEL        = "CANCEL"
	ACTION_CANCEL_REJECT = "CANCEL_REJECT"
	ACTION_EXECUTE       = "EXECUTION"
)

// Action represents an action event passed by the exchange
type Action struct {
	action_type ActionType
	order       Order // Used to represent an action performed on the incoming order
	other_order Order // Used to represent an action performed on the existing book order
	fill_size   Size  // Number of shares filled in the execution
}

// newOrderAction creates a new order action based on the order side (Bid or Ask)
func newOrderAction(order *Order) *Action {
	if order.side == Bid {
		return &Action{
			action_type: ACTION_BID,
			order:       *order,
		}
	} else {
		return &Action{
			action_type: ACTION_ASK,
			order:       *order,
		}
	}
}

// newOrderRejectAction creates a new order rejection action
// This is used in cases of the incoming failing validation (eg. order.price > MAX_PRICE)
func newOrderRejectAction() *Action {
	return &Action{
		action_type: ACTION_ORDER_REJECT,
	}
}

// newCancelAction creates a new cancel action, based on the order to be cancelled
func newCancelAction(order *Order) *Action {
	return &Action{
		action_type: ACTION_CANCEL,
		order:       *order,
	}
}

// newCancelRejectAction creates a new cancel rejection action
// This is used in cases of the cancel OrderID not being found, so the cancel is rejected
func newCancelRejectAction() *Action {
	return &Action{
		action_type: ACTION_CANCEL_REJECT,
	}
}

// TODO: RESTART DOCUMENTATION HERE ****************************************************

// Execution occurs at entry.price for 'price improvement'
func newExecuteAction(order *Order, entry *Order, fill_size Size) *Action {
	if order.side == Bid {
		return &Action{
			action_type: ACTION_EXECUTE,
			order:       *order,
			other_order: *entry,
			fill_size:   fill_size,
		}
	} else {
		return &Action{
			action_type: ACTION_EXECUTE,
			order:       *entry,
			other_order: *order,
			fill_size:   fill_size,
		}
	}
}

// TODO: Unclear if this is actually correctly reporting the side consistently? Check again!!

// String returns a string representation of the action, used for logging
func (action *Action) String() string {
	switch action.action_type {
	case ACTION_BID:
		return fmt.Sprintf(
			"ORDER. ID: %v, Symbol: %v, Side: %v, Price: %v, Size: %v, Trader: %v",
			action.order.order_id,
			action.order.symbol,
			"Bid",
			action.order.price,
			action.order.size,
			action.order.trader,
		)

	case ACTION_ASK:
		return fmt.Sprintf(
			"ORDER. ID: %v, Symbol: %v, Side: %v, Price: %v, Size: %v, Trader: %v",
			action.order.order_id,
			action.order.symbol,
			"Ask",
			action.order.price,
			action.order.size,
			action.order.trader,
		)

	case ACTION_ORDER_REJECT:
		return "ORDER REJECTED"

	case ACTION_CANCEL:
		return fmt.Sprintf("CANCEL. ID: %v", action.order.order_id)

	case ACTION_CANCEL_REJECT:
		return "CANCEL REJECTED"

	case ACTION_EXECUTE:
		return fmt.Sprintf(
			"EXECUTION. Bid_ID: %v, Ask_ID: %v, Symbol: %v, Price: %v, Size: %v, Bid_Trader: %v, Ask_Trader: %v",
			action.order.order_id,
			action.other_order.order_id,
			action.order.symbol,
			action.order.price, // TODO: This is not consistently entry.price
			action.fill_size,
			action.order.trader,
			action.other_order.trader,
		)

	default:
		return fmt.Sprintf("Unknown Action Type: %s", action.action_type)
	}
}
