package main

import (
	"testing"
	"time"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func Test_parseOperation(t *testing.T) {
	tests := []struct {
		name          string
		givenJSON     string
		wantOperation Operation
	}{
		{
			name:      "should parse account operation",
			givenJSON: `{"account": {"active-card": true, "available-limit": 100}}`,
			wantOperation: Operation{
				Account: domain.Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
			},
		},
		{
			name:      "should parse transaction operation",
			givenJSON: `{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			wantOperation: Operation{
				Transaction: domain.Transaction{
					Amount:    20,
					Merchant:  "Burger King",
					CreatedAt: time.Date(2019, 02, 13, 10, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotOperation := parseOperation(test.givenJSON)

			assert.Equal(t, test.wantOperation, gotOperation)
		})
	}
}
