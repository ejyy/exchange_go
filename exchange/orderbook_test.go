package exchange

import (
	"testing"
)

func TestOrderBookInit(t *testing.T) {
	exchange := &Exchange{}
	ob := &OrderBook{}
	ob.init("TEST", exchange)

	if ob.symbol != "TEST" {
		t.Errorf("Expected symbol to be 'TEST', got %s", ob.symbol)
	}
	if ob.exchange != exchange {
		t.Errorf("Expected exchange to be set correctly")
	}
	if ob.asks == nil || ob.bids == nil {
		t.Errorf("Expected btrees to be initialized")
	}
}

func TestOrderBookLimitHandle(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Bid, trader: 1}
	ob.limitHandle(order)

	if len(exchange_engine.actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(exchange_engine.actions))
	}
	if ob.bids.Len() != 1 {
		t.Errorf("Expected 1 bid order in the order book, got %d", ob.bids.Len())
	}
}

func TestOrderBookFillAskSide(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Bid, trader: 1}
	ob.insertIntoBook(&order)

	incomingOrder := Order{order_id: 2, price: 100, size: 5, side: Ask, trader: 2}
	ob.fillAskSide(&incomingOrder)

	if incomingOrder.size != 0 {
		t.Errorf("Expected incoming order to be fully filled, remaining size %d", incomingOrder.size)
	}
	if ob.asks.Len() != 0 {
		t.Errorf("Expected no ask orders in the order book, got %d", ob.asks.Len())
	}
}

func TestOrderBookFillBidSide(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Ask, trader: 1}
	ob.insertIntoBook(&order)

	incomingOrder := Order{order_id: 2, price: 100, size: 5, side: Bid, trader: 2}
	ob.fillBidSide(&incomingOrder)

	if incomingOrder.size != 0 {
		t.Errorf("Expected incoming order to be fully filled, remaining size %d", incomingOrder.size)
	}
	if ob.bids.Len() != 0 {
		t.Errorf("Expected no bid orders in the order book, got %d", ob.bids.Len())
	}
}

func TestOrderBookInsertIntoBook(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Bid, trader: 1}
	ob.insertIntoBook(&order)

	if ob.bids.Len() != 1 {
		t.Errorf("Expected 1 bid order in the order book, got %d", ob.bids.Len())
	}
	if ob.exchange.order_id_map[1] != order {
		t.Errorf("Expected order_id_map to contain the order")
	}
}
