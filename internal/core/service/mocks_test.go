package service

import (
	"time"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type accountRepositoryMock struct {
	mock.Mock
}

func (mock *accountRepositoryMock) SaveAccount(account domain.Account) (domain.Account, error) {
	args := mock.Called(account)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (mock *accountRepositoryMock) FindAccount() (domain.Account, error) {
	args := mock.Called()
	return args.Get(0).(domain.Account), args.Error(1)
}

func (mock *accountRepositoryMock) UpdateAccountLimit(newAvailableLimit int) {
	mock.Called(newAvailableLimit)
}

type transactionRepositoryMock struct {
	mock.Mock
}

func (mock *transactionRepositoryMock) SaveTransaction(transaction domain.Transaction) {
	mock.Called(transaction)
}

func (mock *transactionRepositoryMock) FindTransactionsAfter(time time.Time) []domain.Transaction {
	args := mock.Called(time)
	return args.Get(0).([]domain.Transaction)
}

type accountServicerMock struct {
	mock.Mock
}

func (mock *accountServicerMock) GetAccount() (domain.Account, error) {
	args := mock.Called()
	return args.Get(0).(domain.Account), args.Error(1)
}

func (mock *accountServicerMock) SetAccountLimit(newAvailableLimit int) domain.Account {
	args := mock.Called(newAvailableLimit)
	return args.Get(0).(domain.Account)
}
