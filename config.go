package main

const (
	MAX_PRICE   Price = 100_000
	MIN_PRICE   Price = 1
	EST_ORDERS  Size  = 1_000_000 // Rough estimate of number of orders (to pre-allocate order_id_map)
	EST_SYMBOLS Size  = 1000      // Rough estimate of number of symbols (to pre-allocate orderbooks_map)
)
