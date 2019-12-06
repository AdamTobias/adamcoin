package block

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateBlock(t *testing.T) {
	blk := block{
		Header: header{
			Nonce:   "6046848591118301370",
			PrevBlk: "5678",
			Hash:    []byte{},
		},
		Body: body{
			Txns: []transaction{},
		},
	}
	tests := []struct {
		desc     string
		badHash  bool
		badNonce bool
		err      error
	}{
		{
			desc: "Happy case",
		},
		{
			desc:    "bad hash",
			badHash: true,
			err:     errors.New("bad hash"),
		},
		{
			desc:     "bad nonce",
			badNonce: true,
			err:      errors.New("bad nonce"),
		},
	}
	for _, tt := range tests {
		if tt.badNonce {
			blk.Header.Nonce = "bad nonce"
		}
		newHash := crypto.SHA256
		msg, _ := json.Marshal(blk)
		pssh := newHash.New()
		pssh.Write(msg)
		hashed := pssh.Sum(nil)
		blk.Header.Hash = hashed
		if tt.badHash {
			blk.Header.Hash = []byte("bad hash")
		}

		err := ValidateBlock(blk, unspentCoins{})
		if tt.err != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestValidateTransaction(t *testing.T) {
	pvKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &pvKey.PublicKey
	name, _ := json.Marshal(pubKey)
	nameStr := string(name)
	tests := []struct {
		desc       string
		txn        transaction
		us         unspentCoins
		invalidSig bool
		err        error
	}{
		{
			desc: "Happy case",
			txn: transaction{
				Sender: pubKey,
				Receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				Amount: 50,
			},
			us:  unspentCoins{nameStr: 100},
			err: nil,
		},
		{
			desc: "Amount in does not equal Amount out",
			txn: transaction{
				Sender: pubKey,
				Receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				Amount: 60,
			},
			us:  unspentCoins{nameStr: 100},
			err: errors.New("Amount in does not equal Amount out"),
		},
		{
			desc: "Sender doesn't have enough coins",
			txn: transaction{
				Sender: pubKey,
				Receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				Amount: 60,
			},
			us:  unspentCoins{nameStr: 10},
			err: errors.New("not enough coins"),
		},
		{
			desc: "invalid sig",
			txn: transaction{
				Sender: pubKey,
				Receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				Amount: 50,
			},
			us:         unspentCoins{nameStr: 100},
			invalidSig: true,
			err:        errors.New("invalid sig"),
		},
	}
	for _, tt := range tests {
		newHash := crypto.SHA256
		msg, _ := json.Marshal(tt.txn)
		pssh := newHash.New()
		if tt.invalidSig {
			msg = append(msg, []byte("incorrect")...)
		}
		pssh.Write(msg)
		hashed := pssh.Sum(nil)

		sig, err := rsa.SignPKCS1v15(rand.Reader, pvKey, crypto.SHA256, hashed)
		if err != nil {
			assert.Fail(t, "test faild with err", err)
		}
		tt.txn.Signature = sig
		err = ValidateTransaction(tt.txn, tt.us)
		if tt.err != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
