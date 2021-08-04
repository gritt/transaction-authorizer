package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unknown/authorizer/internal/core/domain"
)

func Test_parseOperation(t *testing.T) {
	tests := []struct {
		name          string
		givenJSON     string
		wantOperation Input
	}{
		{
			name:      "should parse account operation",
			givenJSON: `{"account": {"active-card": true, "available-limit": 100}}`,
			wantOperation: Input{
				Account: domain.Account{
					ActiveCard:     true,
					AvailableLimit: 100,
				},
			},
		},
		{
			name:      "should parse transaction operation",
			givenJSON: `{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			wantOperation: Input{
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
			gotOperation := parseInput(test.givenJSON)

			assert.Equal(t, test.wantOperation, gotOperation)
		})
	}
}
