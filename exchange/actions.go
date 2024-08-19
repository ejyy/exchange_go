package exchange

import "fmt"

// ActionType represents the type of action event passed by the exchange
type ActionType uint8

// Define the action types used in the exchange for various order event states
const (
	ActionBid ActionType = iota
	ActionAsk
	ActionOrderReject
	ActionCancel
	ActionCancelReject
	ActionExecute
)

// Action represents an action event passed by the exchange
type Action struct {
	action_type ActionType
	order       Order // Used to represent an action performed on the incoming order
	cross_order Order // Used to represent an action performed on the existing book order
	fill_size   Size  // Number of shares filled in the execution
	fill_price  Price // Price at which the execution occurrs
}

// newOrderAction creates a new order action based on the order side (Bid or Ask)
func newOrderAction(order *Order) *Action {
	if order.side == Bid {
		return &Action{
			action_type: ActionBid,
			order:       *order,
		}
	} else {
		return &Action{
			action_type: ActionAsk,
			order:       *order,
		}
	}
}

// newOrderRejectAction creates a new order rejection action
// This is used in cases of the incoming failing validation (eg. order.price > MAX_PRICE)
func newOrderRejectAction() *Action {
	return &Action{
		action_type: ActionOrderReject,
	}
}

// newCancelAction creates a new cancel action, based on the order to be cancelled
func newCancelAction(order *Order) *Action {
	return &Action{
		action_type: ActionCancel,
		order:       *order,
	}
}

// newCancelRejectAction creates a new cancel rejection action
// This is used in cases of the cancel OrderID not being found, so the cancel is rejected
func newCancelRejectAction() *Action {
	return &Action{
		action_type: ActionCancelReject,
	}
}

// newExecuteAction creates a new execution action, based on the two orders being executed
// The fill_size is the number of shares filled in the execution
// Execution occurs at entry.price for 'price improvement'
func newExecuteAction(order *Order, entry *Order, fill_size Size) *Action {
	if order.side == Bid {
		return &Action{
			action_type: ActionExecute,
			order:       *order,
			cross_order: *entry,
			fill_size:   fill_size,
			fill_price:  entry.price,
		}
	} else {
		return &Action{
			action_type: ActionExecute,
			order:       *entry,
			cross_order: *order,
			fill_size:   fill_size,
			fill_price:  entry.price,
		}
	}
}

// String returns a string representation of the action, used for logging
func (action *Action) String() string {
	switch action.action_type {
	// String reporting for a new Bid order
	case ActionBid:
		return fmt.Sprintf(
			"ORDER. ID: %v, Symbol: %v, Side: %v, Price: %v, Size: %v, Trader: %v",
			action.order.orderID,
			action.order.symbol,
			"Bid",
			action.order.price,
			action.order.size,
			action.order.trader,
		)

	// String reporting for a new Ask order
	case ActionAsk:
		return fmt.Sprintf(
			"ORDER. ID: %v, Symbol: %v, Side: %v, Price: %v, Size: %v, Trader: %v",
			action.order.orderID,
			action.order.symbol,
			"Ask",
			action.order.price,
			action.order.size,
			action.order.trader,
		)

	// String reporting for an order rejection
	case ActionOrderReject:
		return "ORDER REJECTED"

	// String reporting for a cancel action
	case ActionCancel:
		return fmt.Sprintf("CANCEL. ID: %v", action.order.orderID)

	// String reporting for a cancel rejection
	case ActionCancelReject:
		return "CANCEL REJECTED"

	// String reporting for an execution action
	case ActionExecute:
		// The Bid order is always reported first in the execution action
		return fmt.Sprintf(
			"EXECUTION. Bid_ID: %v, Ask_ID: %v, Symbol: %v, Price: %v, Size: %v, Bid_Trader: %v, Ask_Trader: %v",
			action.order.orderID,       // Bid orderID
			action.cross_order.orderID, // Ask orderID
			action.order.symbol,
			action.fill_price, // As above, execution occurs at entry.price for 'price improvement'
			action.fill_size,
			action.order.trader,       // Bid trader
			action.cross_order.trader, // Ask trader
		)

	// Default case for unknown action types
	default:
		return fmt.Sprintf("Unknown Action Type: %v", action.action_type)
	}
}
