package exchange

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if MaxPrice != 100_000 {
		t.Errorf("Expected MAX_PRICE to be 100000, got %d", MaxPrice)
	}
	if MinPrice != 1 {
		t.Errorf("Expected MIN_PRICE to be 1, got %d", MinPrice)
	}
	if EstNumOrders != 1_000_000 {
		t.Errorf("Expected EST_ORDERS to be 1000000, got %d", EstNumOrders)
	}
	if EstNumSymbols != 1_000 {
		t.Errorf("Expected EST_SYMBOLS to be 1000, got %d", EstNumSymbols)
	}
	if ChanSize != 10_000 {
		t.Errorf("Expected CHAN_SIZE to be 10000, got %d", ChanSize)
	}
}
