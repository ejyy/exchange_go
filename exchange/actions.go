package exchange

import "fmt"

type ActionType string

const (
	ACTION_BID           = "ORDER_BID"
	ACTION_ASK           = "ORDER_ASK"
	ACTION_ORDER_REJECT  = "ORDER_REJECT"
	ACTION_CANCEL        = "CANCEL"
	ACTION_CANCEL_REJECT = "CANCEL_REJECT"
	ACTION_EXECUTE       = "EXECUTION"
)

type Action struct {
	action_type ActionType
	order       Order
	other_order Order
	fill_size   Size
}

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

func newOrderRejectAction() *Action {
	return &Action{
		action_type: ACTION_ORDER_REJECT,
	}
}

func newCancelAction(order *Order) *Action {
	return &Action{
		action_type: ACTION_CANCEL,
		order:       *order,
	}
}

func newCancelRejectAction() *Action {
	return &Action{
		action_type: ACTION_CANCEL_REJECT,
	}
}

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
