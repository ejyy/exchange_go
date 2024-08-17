# exchange_go
Prototype exchange, implementing a multi-symbol limit order book matching engine in Go

## Project todo list
- [x] Initial prototype of limit order book
- [x] Multi-symbol limit order book
- [x] Actions channel here to handle message passing for order and execution reporting
- [x] Order validation within Exchange Limit function
- [x] Pre-warm the Exchange by initialising with a list of symbols
- [x] Public vs private functions based on capitalisation
- [ ] Partial struct storage within order_id_map (just need size and trader)
- [x] Tidy up execution reporting to minimise code repetition
- [x] Move current_order_id and order_id_map out of OrderBook and into Exchange to prevent having to re-make a map for each symbol
- [ ] Thread safety with atomic / mutex
- [ ] TCP / websockets server for order handling and communication
- [ ] README documentation
- [ ] Code documentation
- [ ] Testing suite
