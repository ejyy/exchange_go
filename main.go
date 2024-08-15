package main

func main() {

	var exchange Exchange
	exchange.Init("Example exchange")

	test_order := Order{symbol: "AAPL", price: 100, size: 1000, side: Bid, trader: 1}
	exchange.Limit(test_order)

	test_order2 := Order{symbol: "AAPL", price: 100, size: 1000, side: Ask, trader: 2}
	exchange.Limit(test_order2)

	test_order3 := Order{symbol: "GOOGL", price: 100, size: 1000, side: Bid, trader: 1}
	exchange.Limit(test_order3)

	test_order4 := Order{symbol: "GOOGL", price: 100, size: 1000, side: Ask, trader: 2}
	exchange.Limit(test_order4)

	// exchange.Cancel("APPL", 1)
	// exchange.Cancel("APPL", 2)
}
