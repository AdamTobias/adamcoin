package client

import "blockchain/block"

type Client struct {
	peers map[string]bool
}

func (c *Client) AddPeer(p string) {
	if c.peers == nil {
		c.peers = make(map[string]bool)
	}
	c.peers[p] = true
}

func (c Client) PostBlock(blk block.Block) error {
	for peer, _ := range c.peers {
		// post block to peer
		_ = peer
	}
	return nil
}
