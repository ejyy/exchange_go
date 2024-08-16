package exchange

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if MAX_PRICE != 100_000 {
		t.Errorf("Expected MAX_PRICE to be 100000, got %d", MAX_PRICE)
	}
	if MIN_PRICE != 1 {
		t.Errorf("Expected MIN_PRICE to be 1, got %d", MIN_PRICE)
	}
	if EST_ORDERS != 1_000_000 {
		t.Errorf("Expected EST_ORDERS to be 1000000, got %d", EST_ORDERS)
	}
	if EST_SYMBOLS != 1_000 {
		t.Errorf("Expected EST_SYMBOLS to be 1000, got %d", EST_SYMBOLS)
	}
	if CHAN_SIZE != 10_000 {
		t.Errorf("Expected CHAN_SIZE to be 10000, got %d", CHAN_SIZE)
	}
}
