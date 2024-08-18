package main

import (
	"fmt"

	"github.com/ejyy/exchange_go/exchange"
)

func main() {

	// Create an exchange engine
	var exchange_engine exchange.Exchange

	// Create channel to receive actions from the exchange engine (and a done channel)
	var actions = make(chan *exchange.Action, exchange.CHAN_SIZE)
	var done_channel = make(chan bool)

	// Initialize the exchange engine
	exchange_engine.Init("Example exchange", actions)

	// Pre-warm the exchange engine with some example symbols
	var warming_symbols = []string{"AAPL", "GOOGL"}
	exchange_engine.PreWarmWithSymbols(warming_symbols)

	// Start a goroutine to listen for actions from the exchange engine
	go func() {
		for {
			select {
			case action := <-actions:
				fmt.Printf("Action: %+v\n", action)
			case <-done_channel:
				close(actions)
				return
			}
		}
	}()

	// Send some example orders to the exchange engine
	exchange_engine.Limit("AAPL", 100, 1000, exchange.Bid, 1)
	exchange_engine.Limit("AAPL", 100, 1000, exchange.Ask, 2)
	exchange_engine.Limit("GOOGL", 100, 1000, exchange.Bid, 3)
	exchange_engine.Limit("GOOGL", 100, 1000, exchange.Ask, 4)

	// Send some example cancels to the exchange engine
	exchange_engine.Cancel(1)
	exchange_engine.Cancel(2)

	// Send a done signal to the exchange engine
	// Note, this will close the actions channel meaning producing a variable number of returned messages
	done_channel <- true
}
