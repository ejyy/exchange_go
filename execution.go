package main

type Execution struct {
	symbol       string
	order_id_bid OrderID
	order_id_ask OrderID
	price        Price
	size         Size
	trader_bid   TraderID
	trader_ask   TraderID
}

func (e *Execution) String() string {
	return fmt.Sprintf("EXECUTION... Symbol: %v, Bid ID: %v, Ask ID: %v, Price: %v, Size: %v, Bid trader: %v, Ask trader: %v",
		e.symbol, e.order_id_bid, e.order_id_ask, e.price, e.size, e.trader_bid, e.trader_ask)
}