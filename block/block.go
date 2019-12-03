package block

import (
	"errors"
	"crypto"
	"crypto/rsa"
	"encoding/json"
)

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
	sender *rsa.PublicKey
	receivers unspentCoins
	amount int
	signature []byte
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
	senderName, _ := json.Marshal(txn.sender)
	holding, ok := unspent[string(senderName)]
	if !ok || holding < txn.amount {
		return errors.New("sender does not have enough coins")
	}
	// validate signature
	
	msgObj := transaction{
		sender: txn.sender,
		receivers: txn.receivers,
		amount: txn.amount,
	}
	msg, _  := json.Marshal(msgObj)
	newHash := crypto.SHA256
	pssh := newHash.New()
	pssh.Write(msg)
	hashed := pssh.Sum(nil)

	err := rsa.VerifyPKCS1v15(txn.sender, newHash, hashed, txn.signature); if err != nil {
		return errors.New("invalid sig")
	}
	// txn is valid
	return nil
}
	
