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
	if ob.ask_min != MAX_PRICE+1 {
		t.Errorf("Expected ask_min to be %d, got %d", MAX_PRICE+1, ob.ask_min)
	}
	if ob.bid_max != MIN_PRICE-1 {
		t.Errorf("Expected bid_max to be %d, got %d", MIN_PRICE-1, ob.bid_max)
	}
	for i := range ob.price_points {
		if ob.price_points[i].Len() != 0 {
			t.Errorf("Expected price_points[%d] to be empty", i)
		}
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
	if ob.price_points[100].Len() != 1 {
		t.Errorf("Expected price_points[100] to have 1 order, got %d", ob.price_points[100].Len())
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
	if ob.price_points[100].Len() != 0 {
		t.Errorf("Expected price_points[100] to be empty, got %d", ob.price_points[100].Len())
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
	if ob.price_points[100].Len() != 0 {
		t.Errorf("Expected price_points[100] to be empty, got %d", ob.price_points[100].Len())
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

	if ob.price_points[100].Len() != 1 {
		t.Errorf("Expected price_points[100] to have 1 order, got %d", ob.price_points[100].Len())
	}
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

	if ob.bid_max != 100 {
		t.Errorf("Expected bid_max to be 100, got %d", ob.bid_max)
	}

	order = Order{order_id: 2, price: 50, size: 10, side: Ask, trader: 2}
	ob.insertIntoBook(&order)

	if ob.ask_min != 50 {
		t.Errorf("Expected ask_min to be 50, got %d", ob.ask_min)
	}
}
