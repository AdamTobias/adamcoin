package block

import (
	"crypto"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

const target = "141bc330"

type Block struct {
	Header Header
	Body   body
}

type Header struct {
	Nonce   string
	Hash    []byte
	PrevBlk *Block
}

type body struct {
	Txns []Transaction
}

type Transaction struct {
	Sender    *rsa.PublicKey
	Receivers UnspentCoins
	Amount    int
	Signature []byte
}

type UnspentCoins map[string]int

func ValidateBlock(blk Block, us UnspentCoins) error {
	// validate hash
	msgBlk := Block{
		Body:   blk.Body,
		Header: blk.Header,
	}
	msgBlk.Header.Hash = []byte{}
	msg, _ := json.Marshal(msgBlk)
	newHash := crypto.SHA256
	pssh := newHash.New()
	pssh.Write(msg)
	hashed := pssh.Sum(nil)
	if string(hashed) != string(blk.Header.Hash) {
		return errors.New("incorrect hash")
	}
	// validate hash < target
	valHash := binary.BigEndian.Uint64(blk.Header.Hash)
	valTarget, _ := ParseTarget(target)
	if valHash >= valTarget {
		return errors.New("hash does not meet target criteria")
	}
	// TODO validate prevBlock
	// validate txn list
	for _, txn := range blk.Body.Txns {
		if err := ValidateTransaction(txn, us); err != nil {
			return err
		}
	}
	return nil
}

func ParseTarget(t string) (uint64, error) {
	exp, err := strconv.Atoi(t[:2])
	if err != nil {
		return 0, errors.New("could not parse target -- bad exponent")
	}
	val, err := strconv.ParseUint(target[2:], 16, 32)
	if err != nil {
		return 0, errors.New("could not parse target -- bad value")
	}
	v := uint64(val)
	res := uint64(math.Pow10(exp-3)) * v
	return res, nil
}

func ValidateTransaction(txn Transaction, unspent UnspentCoins) error {
	err := validateAmount(txn)
	if err != nil {
		return err
	}
	err = validateHoldings(txn, unspent)
	if err != nil {
		return err
	}
	err = validateSig(txn)
	if err != nil {
		return err

	}
	// txn is valid
	return nil
}

func validateAmount(txn Transaction) error {
	// validate amount
	totalOut := 0
	for _, amount := range txn.Receivers {
		totalOut = totalOut + amount
	}
	if txn.Amount != totalOut {
		return errors.New("amount in does not equal amount out")
	}
	return nil
}

func validateHoldings(txn Transaction, unspent UnspentCoins) error {
	// validate sender has enough coins
	senderName, _ := json.Marshal(txn.Sender)
	holding, ok := unspent[string(senderName)]
	if !ok || holding < txn.Amount {
		return errors.New("sender does not have enough coins")
	}
	return nil
}

func validateSig(txn Transaction) error {
	// validate signature
	msgObj := Transaction{
		Sender:    txn.Sender,
		Receivers: txn.Receivers,
		Amount:    txn.Amount,
	}
	msg, _ := json.Marshal(msgObj)
	newHash := crypto.SHA256
	pssh := newHash.New()
	pssh.Write(msg)
	hashed := pssh.Sum(nil)

	err := rsa.VerifyPKCS1v15(txn.Sender, newHash, hashed, txn.Signature)
	if err != nil {
		return errors.New("invalid sig")
	}
	return nil
}
