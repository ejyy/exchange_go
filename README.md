# exchange_go
⚡ Exchange_go is a prototype trading exchange, implementing a multi-symbol limit order book matching engine, written in Go. It aims to be high performance, low latency and thread safe. Exchange_go is the order matching engine component of an exchange, taking orders via a function call; it does not currently present a server for client connections.

## Features:
- Multi-symbol limit order book matching engine
- Event reporting via an 'Actions channel' (to handle message passing for order and execution reporting)
- Efficient in-memory model (Btree and Deques for price/time ordering)
- Thread safety (using `sync.mutex`)
- Reasonable test coverage. Reasonable code documentation

> [!WARNING]
> Use in a production environment is strongly discouraged, without much more thorough testing and performance tuning.

## Example usage:
```
package main

import (
	"fmt"

	"github.com/ejyy/exchange_go/exchange"
)

func main() {

	// Create an exchange engine
	var exchange_engine exchange.Exchange

	// Create channel to receive actions from the exchange engine (and a done channel)
	var actions = make(chan *exchange.Action, exchange.ChanSize)
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
```

### Contributions:
Contributions are welcome; please feel free to submit a Pull Request 👍.

## Licence:
MIT
