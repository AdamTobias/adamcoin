package main

import (
	"blockchain/server"
	"fmt"
)

func main() {
	fmt.Println("hello world")
	s := server.Server{}
	s.InitializeServer("1234")
}

type Node struct {
	chain            string
	server           string
	openTransactions []string
}

func (n Node) VerifyTransaction(t transaction) {
	// signature
	// origin transaction(s) is/are valid
}
func (n Node) VerifyBlock(b block) {
	// verify all transactions (do we need to have heard of all the transactions?)
	// verify hash thingy
	// verify the end thingy
}

func (n Node) CreateBlock() {
	// mining function
	// should always be running?

	// look at all openTransactions
	// put a good number of them in a block
	// mine away!
}

type block struct {
	transactionList []transaction
}

type transaction struct{}

type peer struct{}
