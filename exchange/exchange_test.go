package exchange

import (
	"math/rand"
	"testing"
)

func TestExchange_Init(t *testing.T) {
	actions := make(chan *Action, ChanSize)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	if exchange.name != "Test Exchange" {
		t.Errorf("Expected exchange name to be 'Test Exchange', got %s", exchange.name)
	}
	if exchange.currentOrderID != 0 {
		t.Errorf("Expected currentOrderID to be 0, got %d", exchange.currentOrderID)
	}
	if len(exchange.orderbooksMap) != 0 {
		t.Errorf("Expected orderbooksMap to be empty, got %d", len(exchange.orderbooksMap))
	}
	if len(exchange.orderIDMap) != 0 {
		t.Errorf("Expected orderIDMap to be empty, got %d", len(exchange.orderIDMap))
	}
	if exchange.actions != actions {
		t.Errorf("Expected actions channel to be set")
	}
}

func TestExchange_getOrCreateOrderBook(t *testing.T) {
	actions := make(chan *Action, ChanSize)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	symbol := "AAPL"
	orderBook := exchange.getOrCreateOrderBook(symbol)

	if orderBook == nil {
		t.Errorf("Expected order book to be created")
	}
	if exchange.orderbooksMap[symbol] != orderBook {
		t.Errorf("Expected order book to be stored in orderbooksMap")
	}
}

func TestExchange_Limit(t *testing.T) {
	actions := make(chan *Action, ChanSize)
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

func TestExchange_MultipleLimit(t *testing.T) {
	actions := make(chan *Action, ChanSize)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	exchange.Limit("AAPL", 100, 1000, Bid, 1)
	exchange.Limit("AAPL", 100, 1000, Ask, 2)
	exchange.Limit("GOOGL", 200, 10, Bid, 3)
	exchange.Limit("GOOGL", 200, 10, Ask, 4)

	orderBook_aapl := exchange.getOrCreateOrderBook("AAPL")
	if orderBook_aapl.bids.Len() != 0 {
		t.Errorf("Expected AAPL bids book to be empty after full fill")
	}
	if orderBook_aapl.asks.Len() != 0 {
		t.Errorf("Expected AAPL asks book to be empty after full fill")
	}

	orderBook_googl := exchange.getOrCreateOrderBook("GOOGL")
	if orderBook_googl.bids.Len() != 0 {
		t.Errorf("Expected GOOGL bids book to be empty after full fill")
	}
	if orderBook_googl.asks.Len() != 0 {
		t.Errorf("Expected GOOGL asks book to be empty after full fill")
	}
}

// Expand test suite to include order validation tests

func TestExchange_Cancel(t *testing.T) {
	actions := make(chan *Action, ChanSize)
	var exchange Exchange
	exchange.Init("Test Exchange", actions)

	symbol := "AAPL"
	price := Price(100)
	size := Size(10)
	side := Bid
	trader := TraderID(1)

	exchange.Limit(symbol, price, size, side, trader)
	orderID := exchange.currentOrderID

	exchange.Cancel(orderID)

	if exchange.orderIDMap[orderID].size != 0 {
		t.Errorf("Expected order size to be 0 after cancellation")
	}
}

// go test -bench=BenchmarkExchange
func BenchmarkExchange(b *testing.B) {
	minSize := 1
	maxSize := 20
	minPrice := 8000
	maxPrice := 9500

	var actions = make(chan *Action, ChanSize)

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
