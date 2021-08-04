package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unknown/authorizer/internal/core/domain"
)

func TestAuthorizeTransaction(t *testing.T) {
	givenInactiveAccount := domain.Account{
		ActiveCard:     false,
		AvailableLimit: 100,
	}
	givenActiveAccount := domain.Account{
		ActiveCard:     true,
		AvailableLimit: 100,
	}
	givenTransaction := domain.Transaction{
		Amount:    25,
		Merchant:  "ifood",
		CreatedAt: time.Now().UTC(),
	}
	givenTime := givenTransaction.CreatedAt.Add(-2 * time.Minute)

	testCases := map[string]func(*testing.T, *accountServicerMock, *transactionRepositoryMock){
		"should return error when fail to get get account": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			accountServicerMock.On("GetAccount").Return(domain.Account{}, domain.ErrAccountNotInitialized)

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(domain.Transaction{})

			// 	then
			assert.Empty(t, account)
			assert.ElementsMatch(t, errs, []error{domain.ErrAccountNotInitialized})
		},
		"should return error when account card is not active": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			accountServicerMock.On("GetAccount").Return(givenInactiveAccount, nil)

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(domain.Transaction{})

			// 	then
			assert.Equal(t, givenInactiveAccount, account)
			assert.ElementsMatch(t, errs, []error{domain.ErrCardNotActive})
		},
		"should return error when account has insufficient limit": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			accountServicerMock.On("GetAccount").Return(givenActiveAccount, nil)
			transactionRepositoryMock.On("FindTransactionsAfter", mock.AnythingOfType("Time")).Return([]domain.Transaction{})

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(domain.Transaction{Amount: 101})

			// 	then
			assert.Equal(t, givenActiveAccount, account)
			assert.ElementsMatch(t, errs, []error{domain.ErrInsufficientLimit})
		},
		"should return error when high frequency of transactions in small interval": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			foundTransactions := []domain.Transaction{
				{Merchant: "ifood", Amount: 10, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "uber-eats", Amount: 15, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "mercado", Amount: 20, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
			}

			accountServicerMock.On("GetAccount").Return(givenActiveAccount, nil)
			transactionRepositoryMock.On("FindTransactionsAfter", givenTime).Return(foundTransactions)

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(givenTransaction)

			// 	then
			assert.Equal(t, givenActiveAccount, account)
			assert.ElementsMatch(t, errs, []error{domain.ErrHighFrequencySmallInterval})
		},
		"should return error when transactions is doubled in small interval": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			foundTransactions := []domain.Transaction{
				{Merchant: "mercado", Amount: 15, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "ifood", Amount: 25, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
			}

			accountServicerMock.On("GetAccount").Return(givenActiveAccount, nil)
			transactionRepositoryMock.On("FindTransactionsAfter", givenTime).Return(foundTransactions)

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(givenTransaction)

			// 	then
			assert.Equal(t, givenActiveAccount, account)
			assert.ElementsMatch(t, errs, []error{domain.ErrDoubleTransaction})
		},
		"should return list of errors when transaction has multiple violations": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			givenAccount := domain.Account{
				ActiveCard:     true,
				AvailableLimit: 24,
			}
			foundTransactions := []domain.Transaction{
				{Merchant: "mercado", Amount: 15, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "ifood", Amount: 25, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
				{Merchant: "uber-eats", Amount: 100, CreatedAt: time.Now().UTC().Add(-1 * time.Minute)},
			}

			accountServicerMock.On("GetAccount").Return(givenAccount, nil)
			transactionRepositoryMock.On("FindTransactionsAfter", givenTime).Return(foundTransactions)

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(givenTransaction)

			// 	then
			wantErrs := []error{
				domain.ErrDoubleTransaction,
				domain.ErrInsufficientLimit,
				domain.ErrHighFrequencySmallInterval,
			}
			assert.Equal(t, givenAccount, account)
			assert.ElementsMatch(t, errs, wantErrs)
		},
		"should save authorized transaction and set new account limit": func(t *testing.T, accountServicerMock *accountServicerMock, transactionRepositoryMock *transactionRepositoryMock) {
			// 	given
			givenUpdatedAccount := domain.Account{
				ActiveCard:     true,
				AvailableLimit: 75,
			}

			accountServicerMock.On("GetAccount").Return(givenActiveAccount, nil)
			transactionRepositoryMock.On("FindTransactionsAfter", givenTime).Return([]domain.Transaction{})
			accountServicerMock.On("SetAccountLimit", 75).Return(givenUpdatedAccount)
			transactionRepositoryMock.On("SaveTransaction", givenTransaction).Return()

			transactionService := NewTransactionService(transactionRepositoryMock, accountServicerMock)

			// 	when
			account, errs := transactionService.AuthorizeTransaction(givenTransaction)

			// 	then
			assert.Equal(t, givenUpdatedAccount, account)
			assert.Empty(t, errs)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			accountServicerMock := new(accountServicerMock)
			transactionRepositoryMock := new(transactionRepositoryMock)

			run(t, accountServicerMock, transactionRepositoryMock)

			accountServicerMock.AssertExpectations(t)
			transactionRepositoryMock.AssertExpectations(t)
		})
	}
}
