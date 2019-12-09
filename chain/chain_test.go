package chain

import (
	"blockchain/block"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddBlock(t *testing.T) {
	c := Chain{}
	assert.Equal(t, 0, c.length)
	c.AddBlock(block.Block{Header: block.Header{}})
	assert.Equal(t, 1, c.length)
	c.AddBlock(block.Block{Header: block.Header{}})
	assert.Equal(t, 2, c.length)
}
