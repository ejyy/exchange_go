package main

type OrderID uint64
type Side uint8
type Price uint32
type Size uint32
type TraderID uint16

const (
	Bid Side = iota
	Ask
)

type Order struct {
	order_id OrderID
	price    Price // Price in ticks (eg. 12345 would be 123.45)
	size     Size
	side     Side // Bid or Ask
	trader   TraderID
	symbol   string
}