package block

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateTransaction(t *testing.T) {
	tests := []struct{
		desc string
		txn transaction
		us unspentCoins
		err error
	}{
		{
			desc: "Happy case",
			txn : transaction{
				sender: "p0",
		receivers: unspentCoins{
				"p1": 20,
				"p2": 30,
			},
		amount: 50,
	},
	us: unspentCoins{"p0":100},
	err: nil,
},
{
	desc: "amount in does not equal amount out",
	txn : transaction{
		sender: "p0",
		receivers: unspentCoins{
			"p1": 20,
			"p2": 30,
		},
		amount: 60,
	},
	us: unspentCoins{"p0": 100},
	err: errors.New("amount in does not equal amount out"),
},
{
	desc: "sender doesn't have enough coins",
	txn : transaction{
		sender: "p0",
		receivers: unspentCoins{
			"p1": 20,
			"p2": 30,
		},
		amount: 60,
	},
	us: unspentCoins{"p0": 10},
	err: errors.New("not enough coins"),
},
}
for _, tt := range tests {
	err := ValidateTransaction(tt.txn, tt.us) 
	if tt.err != nil {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
	}
}
}
