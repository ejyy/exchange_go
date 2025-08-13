package main

import (
	"github.com/ejyy/exchange_go/exchange"
	"github.com/ejyy/exchange_go/server"
)

func main() {

	// Create an exchange engine
	var exchange_engine exchange.Exchange

	// Create channel to receive actions from the exchange engine
	var actions = make(chan *exchange.Action, exchange.ChanSize)

	// Initialize the exchange engine
	exchange_engine.Init("Example exchange", actions)

	// Pre-warm the exchange engine with some example symbols
	var warming_symbols = []string{"AAPL", "GOOGL"}
	exchange_engine.PreWarmWithSymbols(warming_symbols)

	// Create a server engine
	var server server.Server

	// Initialize the server engine
	server.Init(&exchange_engine, actions)
}
