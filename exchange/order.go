package exchange

// Define the types used in the exchange. These are used to represent the orderbook, orders, and traders
type OrderID uint64  // Unique identifier for an order [range 0-2^64]
type Side uint8      // Bid or Ask
type Price uint32    // Price in ticks (eg. 12345 would be 123.45) [range 0-2^32]
type Size uint32     // Size integer [range 0-2^32]
type TraderID uint16 // Unique identifier for a trader [range 0-2^16]

// Define the two sides of an order
const (
	Bid Side = iota // Bid side represents a buy order
	Ask             // Ask side represents a sell order
)

// Order represents an order on the exchange
type Order struct {
	order_id OrderID
	price    Price
	size     Size
	side     Side
	trader   TraderID
	symbol   string // Symbol of the order (eg. AAPL, GOOGL)
}
