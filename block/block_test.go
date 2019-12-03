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
		err error
	}{
		{
			desc: "Happy case",
			txn : transaction{
		receivers: []receiver{
			{
				name: "p1",
				amount: 20,
			}, {
				name: "p2",
				amount: 30,
			},
		},
		amount: 50,
	},
	err: nil,
},
{
	desc: "amount in does not equal amount out",
	txn : transaction{
		receivers: []receiver{
			{
				name: "p1",
				amount: 20,
			}, {
				name: "p2",
				amount: 30,
			},
		},
		amount: 60,
	},
	err: errors.New("amount in does not equal amount out"),
},
}
for _, tt := range tests {
	err := ValidateTransaction(tt.txn)
	if tt.err != nil {
		assert.Error(t, err)
	} else {
		assert.Nil(t, err)
	}
}
}
