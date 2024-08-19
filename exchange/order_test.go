package exchange

import (
	"testing"
)

func TestOrderCreation(t *testing.T) {
	order := Order{
		orderID: 1,
		price:   12345,
		size:    100,
		side:    Bid,
		trader:  1,
		symbol:  "TEST",
	}

	if order.orderID != 1 {
		t.Errorf("Expected orderID to be 1, got %d", order.orderID)
	}
	if order.price != 12345 {
		t.Errorf("Expected price to be 12345, got %d", order.price)
	}
	if order.size != 100 {
		t.Errorf("Expected size to be 100, got %d", order.size)
	}
	if order.side != Bid {
		t.Errorf("Expected side to be Bid, got %d", order.side)
	}
	if order.trader != 1 {
		t.Errorf("Expected trader to be 1, got %d", order.trader)
	}
	if order.symbol != "TEST" {
		t.Errorf("Expected symbol to be 'TEST', got %s", order.symbol)
	}
}

func TestOrderSide(t *testing.T) {
	bidOrder := Order{side: Bid}
	askOrder := Order{side: Ask}

	if bidOrder.side != Bid {
		t.Errorf("Expected side to be Bid, got %d", bidOrder.side)
	}
	if askOrder.side != Ask {
		t.Errorf("Expected side to be Ask, got %d", askOrder.side)
	}
}

func TestOrderPrice(t *testing.T) {
	order := Order{price: 12345}

	if order.price != 12345 {
		t.Errorf("Expected price to be 12345, got %d", order.price)
	}
}

func TestOrderSize(t *testing.T) {
	order := Order{size: 100}

	if order.size != 100 {
		t.Errorf("Expected size to be 100, got %d", order.size)
	}
}
