package controller

import (
	"blockchain/block"
	"blockchain/chain"
	"blockchain/controller/client"
	"blockchain/controller/miner"
	"encoding/json"
)

type Controller struct {
	blockchain chain.Chain
	coinPool   block.UnspentCoins
	cli        client.Client
	mnr        miner.Miner
	looseTxns  []block.Transaction
}

func NewController(peer string) Controller {
	cli := client.Client{}
	cli.AddPeer(peer)
	// get chain from peer
	// init chain
	return Controller{cli: cli}
}

func (c Controller) AddPeer(p string) {
	c.cli.AddPeer(p)
}

func (c *Controller) AddTxn(t block.Transaction) {
	// use map to know if we have already seen this txn?
	c.looseTxns = append(c.looseTxns, t)
	// update miner -- if new txn can go in block, need to stop mining and mine with additional txn.  or maybe mining only occurs with exactly X blocks so this wont be necessary?
}

func (c *Controller) AddBlock(blk block.Block) error {
	// validate block
	if err := block.ValidateBlock(blk, c.coinPool); err != nil {
		return err
	}
	// add block
	c.blockchain.AddBlock(blk)
	// update coinPool
	for _, txn := range blk.Body.Txns {
		for name, amount := range txn.Receivers {
			if _, ok := c.coinPool[name]; !ok {
				c.coinPool[name] = 0
			}
			c.coinPool[name] += amount
		}
		senderName, _ := json.Marshal(txn.Sender)
		c.coinPool[string(senderName)] -= txn.Amount
	}
	// update mining
	// post block to peers
	return nil
}
