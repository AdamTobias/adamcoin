package chain

import (
	"blockchain/block"
)

type Chain struct {
	head   block.Block
	length int
}

func (c *Chain) AddBlock(blk block.Block) error {
	blk.Header.PrevBlk = &c.head
	c.head = blk
	c.length++
	return nil
}
