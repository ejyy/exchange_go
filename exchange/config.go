package exchange

// Define the constants used in the exchange
// These are used to provide bounds for the exchange and pre-allocate memory
const (
	MAX_PRICE   Price = 100_000
	MIN_PRICE   Price = 1
	EST_ORDERS  Size  = 1_000_000 // Rough estimate of number of orders (to pre-allocate order_id_map)
	EST_SYMBOLS Size  = 1_000     // Rough estimate of number of symbols (to pre-allocate orderbooks_map)
	CHAN_SIZE   Size  = 10_000    // Channel buffer size
)
