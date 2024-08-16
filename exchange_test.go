package main

import (
	"math/rand"
	"testing"
)

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
		for range exchange.actions {
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

		exchange.Limit(Order{symbol: symbol, price: Price(price), size: Size(size), side: side, trader: 1})
	}
}
