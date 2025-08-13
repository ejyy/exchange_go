package main

import (
	"fmt"

	"github.com/ejyy/exchange_go/exchange"
	"github.com/ejyy/exchange_go/server"
)

func main() {

	// Create an exchange engine
	var exchange_engine exchange.Exchange

	// Create channel to receive actions from the exchange engine (and a done channel)
	var actions = make(chan *exchange.Action, exchange.ChanSize)

	// Initialize the exchange engine
	exchange_engine.Init("Example exchange", actions)

	// Pre-warm the exchange engine with some example symbols
	var warming_symbols = []string{"AAPL", "GOOGL"}
	exchange_engine.PreWarmWithSymbols(warming_symbols)

	server := server.NewServer(&exchange_engine, actions)

	fmt.Println("Starting TCP server...")
	if err := server.Start(exchange.TCPPort); err != nil {
		fmt.Println("Server failed: ", err)
	}
}
