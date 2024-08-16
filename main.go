package main

import (
	"fmt"
)

func main() {

	var actions = make(chan *Action, CHAN_SIZE)
	var done_channel = make(chan bool)

	var exchange Exchange
	exchange.Init("Example exchange", actions)

	go func() {
		for {
			select {
			case action := <-exchange.actions:
				fmt.Printf("Action: %+v\n", action)
			case <-done_channel:
				close(actions)
				return
			}
		}
	}()

	test_order := Order{symbol: "AAPL", price: 100, size: 1000, side: Bid, trader: 1}
	exchange.Limit(test_order)

	test_order2 := Order{symbol: "AAPL", price: 100, size: 1000, side: Ask, trader: 2}
	exchange.Limit(test_order2)

	test_order3 := Order{symbol: "GOOGL", price: 100, size: 1000, side: Bid, trader: 3}
	exchange.Limit(test_order3)

	test_order4 := Order{symbol: "GOOGL", price: 100, size: 1000, side: Ask, trader: 4}
	exchange.Limit(test_order4)

	exchange.Cancel(1)
	exchange.Cancel(2)

	done_channel <- true

}
