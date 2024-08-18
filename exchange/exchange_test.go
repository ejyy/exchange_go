package exchange

import (
	"math/rand"
	"testing"
)

func TestExchange_Init(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	if exchange.name != "Test Exchange" {
		t.Errorf("Expected exchange name to be 'Test Exchange', got %s", exchange.name)
	}
	if exchange.current_order_id != 0 {
		t.Errorf("Expected current_order_id to be 0, got %d", exchange.current_order_id)
	}
	if len(exchange.orderbooks_map) != 0 {
		t.Errorf("Expected orderbooks_map to be empty, got %d", len(exchange.orderbooks_map))
	}
	if len(exchange.order_id_map) != 0 {
		t.Errorf("Expected order_id_map to be empty, got %d", len(exchange.order_id_map))
	}
	if exchange.actions != actions {
		t.Errorf("Expected actions channel to be set")
	}
}

func TestExchange_getOrCreateOrderBook(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	symbol := "AAPL"
	orderBook := exchange.getOrCreateOrderBook(symbol)

	if orderBook == nil {
		t.Errorf("Expected order book to be created")
	}
	if exchange.orderbooks_map[symbol] != orderBook {
		t.Errorf("Expected order book to be stored in orderbooks_map")
	}
}

func TestExchange_Limit(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	symbol := "AAPL"
	price := Price(100)
	size := Size(10)
	side := Bid
	trader := TraderID(1)

	exchange.Limit(symbol, price, size, side, trader)

	orderBook := exchange.getOrCreateOrderBook(symbol)
	if orderBook.bids.Len() == 0 {
		t.Errorf("Expected order to be added to the order book")
	}
}

// Expand test suite to include order validation tests

func TestExchange_Cancel(t *testing.T) {
	actions := make(chan *Action, CHAN_SIZE)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	symbol := "AAPL"
	price := Price(100)
	size := Size(10)
	side := Bid
	trader := TraderID(1)

	exchange.Limit(symbol, price, size, side, trader)
	orderID := exchange.current_order_id

	exchange.Cancel(orderID)

	if exchange.order_id_map[orderID].size != 0 {
		t.Errorf("Expected order size to be 0 after cancellation")
	}
}

// go test -bench=BenchmarkExchange
func BenchmarkExchange(b *testing.B) {
	minSize := 1
	maxSize := 20
	minPrice := 8000
	maxPrice := 9500

	var actions = make(chan *Action, CHAN_SIZE)

	var exchange Exchange
	exchange.Init("Test exchange", actions)

	go func() {
		for range actions {
			// Do nothing, just discard the action (avoid overhead of STDOUT)
		}
	}()

	for i := 0; i < b.N; i++ {
		price := rand.Intn(maxPrice-minPrice) + minPrice
		size := rand.Intn(maxSize-minSize) + minSize

		charset := "abcdefghijklmnopqrstuvwxyz"
		symbol := string(charset[rand.Intn(len(charset))])

		var side Side
		if rand.Intn(1000) >= 500 {
			side = Bid
		} else {
			side = Ask
		}

		exchange.Limit(symbol, Price(price), Size(size), side, 1)
	}
}
