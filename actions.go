package main

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

// TODO: String printing function to dereference order pointers for pretty output (avoids the structs filled with 0 when unpopulated)

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
