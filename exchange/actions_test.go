package exchange

import (
	"testing"
)

func TestNewOrderAction(t *testing.T) {
	order := &Order{orderID: 1, symbol: "AAPL", side: Bid, price: 150, size: 10, trader: 1}
	action := newOrderAction(order)
	if action.action_type != ActionBid {
		t.Errorf("Expected action type to be %v, got %v", ActionBid, action.action_type)
	}
	if action.order != *order {
		t.Errorf("Expected order to be %v, got %v", *order, action.order)
	}
}

func TestNewCancelAction(t *testing.T) {
	order := &Order{orderID: 1, symbol: "AAPL", side: Bid, price: 150, size: 0, trader: 1}
	action := newCancelAction(order)
	if action.action_type != ActionCancel {
		t.Errorf("Expected action type to be %v, got %v", ActionCancel, action.action_type)
	}
	if action.order != *order {
		t.Errorf("Expected order to be %v, got %v", *order, action.order)
	}
	if action.order.size != 0 {
		t.Errorf("Expected order to be %v, got %v", 0, action.order.size)
	}
}

func TestNewCancelRejectAction(t *testing.T) {
	action := newCancelRejectAction()
	if action.action_type != ActionCancelReject {
		t.Errorf("Expected action type to be %v, got %v", ActionCancelReject, action.action_type)
	}
}

func TestNewExecuteAction(t *testing.T) {
	order := &Order{orderID: 1, symbol: "AAPL", side: Bid, price: 150, size: 10, trader: 1}
	entry := &Order{orderID: 2, symbol: "AAPL", side: Ask, price: 150, size: 10, trader: 2}
	fill_size := Size(10)
	action := newExecuteAction(order, entry, fill_size)
	if action.action_type != ActionExecute {
		t.Errorf("Expected action type to be %v, got %v", ActionExecute, action.action_type)
	}
	if action.order != *order {
		t.Errorf("Expected order to be %v, got %v", *order, action.order)
	}
	if action.cross_order != *entry {
		t.Errorf("Expected other order to be %v, got %v", *entry, action.cross_order)
	}
	if action.fill_size != fill_size {
		t.Errorf("Expected fill size to be %v, got %v", fill_size, action.fill_size)
	}
}

func TestActionString(t *testing.T) {
	order := &Order{orderID: 1, symbol: "AAPL", side: Bid, price: 150, size: 10, trader: 1}
	entry := &Order{orderID: 2, symbol: "AAPL", side: Ask, price: 150, size: 5, trader: 2}
	fill_size := Size(5)

	tests := []struct {
		action *Action
		want   string
	}{
		{newOrderAction(order), "ORDER. ID: 1, Symbol: AAPL, Side: Bid, Price: 150, Size: 10, Trader: 1"},
		{newCancelAction(order), "CANCEL. ID: 1"},
		{newCancelRejectAction(), "CANCEL REJECTED"},
		{newExecuteAction(order, entry, fill_size), "EXECUTION. Bid_ID: 1, Ask_ID: 2, Symbol: AAPL, Price: 150, Size: 5, Bid_Trader: 1, Ask_Trader: 2"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.action.String(); got != tt.want {
				t.Errorf("Action.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
