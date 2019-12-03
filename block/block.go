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
	receivers []receiver
	amount int
	signature string
}

type receiver struct {
	name string
	amount int
}

func ValidateTransaction(txn transaction) error {
	// validate amount
	totalOut := 0
	for _, rec := range txn.receivers {
		totalOut = totalOut + rec.amount
	}
	if txn.amount != totalOut {
		return errors.New("amount in does not equal amount out")
	}
	return nil
	// validate sender has enough coins
	// validate signature
}
	
