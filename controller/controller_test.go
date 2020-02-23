package controller

import (
	"blockchain/block"
	"testing"
	"time"
)

func TestMine(t *testing.T) {
	ctrl := NewController("test_peer")
	ctrl.Mine(block.Block{})
	time.Sleep(100 * time.Microsecond)
	ctrl.Mine(block.Block{})
	time.Sleep(100 * time.Microsecond)
}
