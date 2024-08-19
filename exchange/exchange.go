package exchange

import (
	"fmt"
	"sync"
)

// Exchange represents the exchange engine, that stores the orderbooks (per symbol) and manages the orders
type Exchange struct {
	name           string
	orderbooksMap  map[string]*OrderBook
	currentOrderID OrderID
	orderIDMap     map[OrderID]Order // TODO: Wasteful to store entire 'Order' struct, only need trader + size (https://go.dev/play/p/4KPix5OEXJC)
	actions        chan *Action
	mutex          sync.RWMutex
}

// TODO: Consider using alternative sync representations here
// Use sync.map instead of a mutex-protected map
// Use atomic operations for the currentOrderID

// Init initialises the exchange with the given name and actions channel, and establishes the order storage
func (ex *Exchange) Init(name string, actions chan *Action) {
	// Lock the exchange mutex to prevent concurrent access
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	ex.name = name
	ex.currentOrderID = 0

	// Pre-allocate the maps to avoid resizing based on estimated values (in config)
	ex.orderbooksMap = make(map[string]*OrderBook, EstNumSymbols)
	ex.orderIDMap = make(map[OrderID]Order, EstNumOrders)

	ex.actions = actions

	// Report the exchange is ready to accept orders via STDOUT
	fmt.Println("Exchange started:", ex.name, "- Ready to accept orders")
}

// getNextOrderID returns the next available order ID in the exchange and increments the counter
func (ex *Exchange) getNextOrderID() OrderID {
	// Lock the exchange mutex to prevent concurrent access
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	ex.currentOrderID += 1
	return ex.currentOrderID
}

// getOrCreateOrderBook returns the orderbook for the given symbol, creating it if it doesn't exist
func (ex *Exchange) getOrCreateOrderBook(symbol string) *OrderBook {
	// Lock the exchange mutex to prevent concurrent access
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	order_book, exists := ex.orderbooksMap[symbol]
	if !exists {
		order_book = new(OrderBook)
		order_book.init(symbol, ex)
		ex.orderbooksMap[symbol] = order_book
	}
	return order_book
}

// PreWarmWithSymbols 'pre-warms' the exchange with the given symbols, creating the orderbooks
// Used to avoid the first order for a symbol being slow due to orderbook creation
func (ex *Exchange) PreWarmWithSymbols(symbols []string) {
	for _, symbol := range symbols {
		ex.getOrCreateOrderBook(symbol)
	}
}

// validateOrder checks the incoming order for validity, ensuring the fields are within bounds
// Used to prevent invalid orders from being processed
func validateOrder(symbol string, price Price, size Size, side Side, trader TraderID) bool {
	if symbol == "" {
		return false
	}
	if price < MinPrice || price > MaxPrice {
		return false
	}
	if size <= 0 {
		return false
	}
	if side != Bid && side != Ask {
		return false
	}
	// TraderID is not specifically validated, as it can be any positive integer
	// This could be extended to check for a valid traderID from a database call or similar
	if trader <= 0 {
		return false
	}
	return true
}

// Limit processes an incoming limit order, validating it and passing it to the appropriate orderbook
func (ex *Exchange) Limit(symbol string, price Price, size Size, side Side, trader TraderID) {
	// Validate the incoming order, rejecting if invalid
	if !validateOrder(symbol, price, size, side, trader) {
		// Report the rejection to the exchange via the actions channel
		ex.actions <- newOrderRejectAction()
		return
	}

	// Initialise the incoming order with the given values
	incomingOrder := Order{
		symbol: symbol,
		price:  price,
		size:   size,
		side:   side,
		trader: trader,
	}

	// Get or create the orderbook for the symbol and process the incoming order
	ob := ex.getOrCreateOrderBook(incomingOrder.symbol)
	incomingOrder.orderID = ex.getNextOrderID()
	ob.limitHandle(incomingOrder)
}

// Cancel processes an incoming cancel order, cancelling the order if it exists in the exchange
func (ex *Exchange) Cancel(orderID OrderID) {
	// Lock the exchange mutex to prevent concurrent access
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	// Check if the order exists in the exchange
	if cancelOrder, ok := ex.orderIDMap[orderID]; ok {
		// If the order size is zero, it has already been cancelled
		if cancelOrder.size == 0 {
			// Report the cancel rejection to the exchange via the actions channel
			ex.actions <- newCancelRejectAction()
		} else {
			// Update the order size to zero to show it has been cancelled
			cancelOrder.size = 0

			// Update the orderIDMap with the cancelled order
			ex.orderIDMap[orderID] = cancelOrder

			// Report the cancellation to the exchange via the actions channel
			ex.actions <- newCancelAction(&cancelOrder)
		}
	} else {
		// If the orderID is not found in the orderIDMap, it cannot be cancelled
		// Report the cancel rejection to the exchange via the actions channel
		ex.actions <- newCancelRejectAction()
	}
}
