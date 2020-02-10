package server

import (
	"blockchain/block"
	"crypto/rsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"testing"
)

func TestParseAddTxnReq(t *testing.T) {
	// happy case
	// incorrect number of rec N E or amt
	// valid rec N (multiple fail cases?)
	// empty sig
	// invalid amt
	// invalid sender
	req, err := http.NewRequest("GET", "http://localhost:9421/addtxn", nil)
	if err != nil {
		fmt.Println("error forming request", err)
	}
	q := req.URL.Query()
	q.Add("senderN", "1234")
	q.Add("senderE", "1234")
	q.Add("recN", "4321")
	q.Add("recE", "65537")
	q.Add("recAmt", "50")
	q.Add("sig", "A12FF42988BC")
	q.Add("amt", "50")
	req.URL.RawQuery = q.Encode()

	s := Server{}
	txn, err := s.ParseAddTxnReq(req)
	pubKey := rsa.PublicKey{
		N: big.NewInt(1234),
		E: 1234,
	}
	exp := block.Transaction{
		Sender: &pubKey,
		Receivers: block.UnspentCoins{
			"{\"N\":4321,\"E\":65537}": 50,
		},
		Amount:    50,
		Signature: []byte("A12FF42988BC"),
	}
	assert.Equal(t, exp, txn)

}
