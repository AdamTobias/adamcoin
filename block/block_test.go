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
				sender: pubKey,
				receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				amount: 50,
			},
			us:  unspentCoins{nameStr: 100},
			err: nil,
		},
		{
			desc: "amount in does not equal amount out",
			txn: transaction{
				sender: pubKey,
				receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				amount: 60,
			},
			us:  unspentCoins{nameStr: 100},
			err: errors.New("amount in does not equal amount out"),
		},
		{
			desc: "sender doesn't have enough coins",
			txn: transaction{
				sender: pubKey,
				receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				amount: 60,
			},
			us:  unspentCoins{nameStr: 10},
			err: errors.New("not enough coins"),
		},
		{
			desc: "invalid sig",
			txn: transaction{
				sender: pubKey,
				receivers: unspentCoins{
					"p1": 20,
					"p2": 30,
				},
				amount: 50,
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
		tt.txn.signature = sig
		err = ValidateTransaction(tt.txn, tt.us)
		if tt.err != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
