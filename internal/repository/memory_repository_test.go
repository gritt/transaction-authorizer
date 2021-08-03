package repository

import (
	"testing"
	"time"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestSaveTransaction(t *testing.T) {
	testCases := map[string]func(*testing.T){
		"should save one transaction with success": func(t *testing.T) {
			// 	given
			givenTransaction := domain.Transaction{
				Amount:    100,
				Merchant:  "ifood",
				CreatedAt: time.Now().UTC(),
			}

			repository := NewMemoryRepository()

			// 	when
			repository.SaveTransaction(givenTransaction)

			// 	then
			wantTransactions := []domain.Transaction{givenTransaction}
			assert.ElementsMatch(t, repository.transactions, wantTransactions)
		},
		"should save multiple transactions with success": func(t *testing.T) {
			// 	given
			givenTransactions := []domain.Transaction{
				{Merchant: "ifood", Amount: 100, CreatedAt: time.Now().UTC()},
				{Merchant: "uber-eats", Amount: 75, CreatedAt: time.Now().UTC()},
				{Merchant: "mercado", Amount: 200, CreatedAt: time.Now().UTC()},
			}

			repository := NewMemoryRepository()

			// 	when
			repository.SaveTransaction(givenTransactions[0])
			repository.SaveTransaction(givenTransactions[1])
			repository.SaveTransaction(givenTransactions[2])

			// 	then
			assert.ElementsMatch(t, repository.transactions, givenTransactions)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestFindTransactionsAfter(t *testing.T) {
	testCases := map[string]func(*testing.T){
		"should return found transactions after given time": func(t *testing.T) {
			// 	given
			givenTransactions := []domain.Transaction{
				{Merchant: "ifood", Amount: 100, CreatedAt: time.Now().UTC().Add(-3 * time.Minute)},
				{Merchant: "uber-eats", Amount: 75, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "mercado", Amount: 200, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
			}

			repository := NewMemoryRepository()

			repository.SaveTransaction(givenTransactions[0])
			repository.SaveTransaction(givenTransactions[1])
			repository.SaveTransaction(givenTransactions[2])

			// 	when
			givenTime := time.Now().UTC().Add(-2 * time.Minute)
			foundTransactions := repository.FindTransactionsAfter(givenTime)

			// 	then
			wantTransactions := []domain.Transaction{
				givenTransactions[1],
				givenTransactions[2],
			}
			assert.ElementsMatch(t, wantTransactions, foundTransactions)
		},
		"should return empty transactions not found after given time": func(t *testing.T) {
			// 	given
			givenTransactions := []domain.Transaction{
				{Merchant: "ifood", Amount: 100, CreatedAt: time.Now().UTC().Add(-3 * time.Minute)},
				{Merchant: "uber-eats", Amount: 75, CreatedAt: time.Now().UTC().Add(-3 * time.Minute)},
				{Merchant: "mercado", Amount: 200, CreatedAt: time.Now().UTC().Add(-3 * time.Minute)},
			}

			repository := NewMemoryRepository()

			repository.SaveTransaction(givenTransactions[0])
			repository.SaveTransaction(givenTransactions[1])
			repository.SaveTransaction(givenTransactions[2])

			// 	when
			givenTime := time.Now().UTC().Add(-2 * time.Minute)
			foundTransactions := repository.FindTransactionsAfter(givenTime)

			// 	then
			assert.Empty(t, foundTransactions)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestSaveAccount(t *testing.T) {
	testCases := map[string]func(*testing.T){
		"should save account with success": func(t *testing.T) {
			// 	given
			givenAccount := domain.Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			}
			repository := NewMemoryRepository()

			// 	when
			savedAccount, err := repository.SaveAccount(givenAccount)

			// 	then
			assert.Equal(t, givenAccount, savedAccount)
			assert.NoError(t, err)
		},
		"should return error when account already initialized": func(t *testing.T) {
			// 	given
			wantAccount := domain.Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			}

			repository := NewMemoryRepository()

			firstAccount, err := repository.SaveAccount(domain.Account{ActiveCard: false, AvailableLimit: 100})
			assert.Equal(t, wantAccount, firstAccount)
			assert.NoError(t, err)

			// 	when
			secondAccount, err := repository.SaveAccount(domain.Account{ActiveCard: true, AvailableLimit: 300})

			// 	then
			assert.Equal(t, wantAccount, secondAccount)
			assert.EqualError(t, err, "account already initialized")
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestFindAccount(t *testing.T) {
	testCases := map[string]func(*testing.T){
		"should find account with success": func(t *testing.T) {
			// 	given
			repository := NewMemoryRepository()
			savedAccount, err := repository.SaveAccount(domain.Account{ActiveCard: false, AvailableLimit: 100})
			assert.NoError(t, err)
			assert.NotEmpty(t, savedAccount)

			// 	when
			foundAccount, err := repository.FindAccount()

			// 	then
			assert.Equal(t, savedAccount, foundAccount)
			assert.NoError(t, err)
		},
		"should return error when account not initialized": func(t *testing.T) {
			// 	given
			repository := NewMemoryRepository()

			// 	when
			foundAccount, err := repository.FindAccount()

			// 	then
			assert.Empty(t, foundAccount)
			assert.EqualError(t, err, "account not initialized")
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestUpdateAccountLimit(t *testing.T) {
	testCases := map[string]func(*testing.T){
		"should update account limit with success": func(t *testing.T) {
			// 	given
			repository := NewMemoryRepository()
			initialAccount, err := repository.SaveAccount(domain.Account{ActiveCard: true, AvailableLimit: 100})
			assert.NoError(t, err)
			assert.NotEmpty(t, initialAccount)

			// 	when
			repository.UpdateAccountLimit(75)

			// 	then
			updatedAccount, err := repository.FindAccount()
			assert.NoError(t, err)
			assert.Equal(t, 75, updatedAccount.AvailableLimit)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
