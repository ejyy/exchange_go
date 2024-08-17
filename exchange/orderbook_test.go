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
}

func TestOrderBookInsertIntoBook(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Bid, trader: 1}
	ob.insertIntoBook(&order)

	if ob.exchange.order_id_map[1] != order {
		t.Errorf("Expected order_id_map to contain the order")
	}
}

func TestOrderBookUpdateBidMaxAskMin(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange_engine Exchange
	exchange_engine.Init("TEST", actions)

	ob := &OrderBook{}
	ob.init("TEST", &exchange_engine)

	order := Order{order_id: 1, price: 100, size: 10, side: Bid, trader: 1}
	ob.insertIntoBook(&order)

	order = Order{order_id: 2, price: 50, size: 10, side: Ask, trader: 2}
	ob.insertIntoBook(&order)
}
