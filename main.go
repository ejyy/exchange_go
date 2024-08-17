package main

import (
	"fmt"

	"github.com/ejyy/exchange_go/exchange"
)

func main() {

	var exchange_engine exchange.Exchange

	var actions = make(chan *exchange.Action, exchange.CHAN_SIZE)
	var done_channel = make(chan bool)

	exchange_engine.Init("Example exchange", actions)

	var warming_symbols = []string{"AAPL", "GOOGL"}
	exchange_engine.PreWarmWithSymbols(warming_symbols)

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

	exchange_engine.Limit("AAPL", 100, 1000, exchange.Bid, 1)
	exchange_engine.Limit("AAPL", 100, 1000, exchange.Ask, 2)
	exchange_engine.Limit("GOOGL", 100, 1000, exchange.Bid, 3)
	exchange_engine.Limit("GOOGL", 100, 1000, exchange.Ask, 4)

	exchange_engine.Cancel(1)
	exchange_engine.Cancel(2)

	done_channel <- true

	// TODO: Getting variable number of messages, based on execution time of below orders. Due to sending done immediately after orders
}
