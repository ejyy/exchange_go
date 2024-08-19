package exchange

// Define the constants used in the exchange
// These are used to provide bounds for the exchange and pre-allocate memory
const (
	MaxPrice      Price = 100_000
	MinPrice      Price = 1
	EstNumOrders  Size  = 1_000_000 // Rough estimate of number of orders (to pre-allocate orderIDMap)
	EstNumSymbols Size  = 1_000     // Rough estimate of number of symbols (to pre-allocate orderbooksMap)
	ChanSize      Size  = 10_000    // Channel buffer size
)
