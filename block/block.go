package block

import "errors"

type block struct {
	header header
	body body
}

type header struct {
	nonce string
	hash string
	prevBlk string
}

type body struct {
	txns []transaction
}

type transaction struct {
	sender string
	receivers unspentCoins
	amount int
	signature string
}

type unspentCoins map[string]int

func ValidateTransaction(txn transaction, unspent unspentCoins) error {
	// validate amount
	totalOut := 0
	for _, amount := range txn.receivers {
		totalOut = totalOut + amount 
	}
	if txn.amount != totalOut {
		return errors.New("amount in does not equal amount out")
	}
	// validate sender has enough coins
	holding, ok := unspent[txn.sender]
	if !ok || holding < txn.amount {
		return errors.New("sender does not have enough coins")
	}
	
	return nil
	// validate signature
}
	
