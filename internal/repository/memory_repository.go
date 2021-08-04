package repository

import (
	"errors"
	"time"

	"github.com/unknown/authorizer/internal/core/domain"
)

type MemoryRepository struct {
	transactions       []domain.Transaction
	account            domain.Account
	accountInitialized bool
}

func NewMemoryRepository() MemoryRepository {
	return MemoryRepository{}
}

func (m *MemoryRepository) SaveTransaction(transaction domain.Transaction) {
	m.transactions = append(m.transactions, transaction)
}

func (m *MemoryRepository) FindTransactionsAfter(time time.Time) []domain.Transaction {
	foundTransactions := []domain.Transaction{}
	for _, transaction := range m.transactions {
		if transaction.CreatedAt.After(time) {
			foundTransactions = append(foundTransactions, transaction)
		}
	}
	return foundTransactions
}

func (m *MemoryRepository) SaveAccount(account domain.Account) (domain.Account, error) {
	if !m.accountInitialized {
		m.account = account
		m.accountInitialized = true
		return m.account, nil
	}
	return m.account, errors.New("account already initialized")
}

func (m *MemoryRepository) FindAccount() (domain.Account, error) {
	if !m.accountInitialized {
		return domain.Account{}, errors.New("account not initialized")
	}
	return m.account, nil
}

func (m *MemoryRepository) UpdateAccountLimit(newAvailableLimit int) {
	m.account.AvailableLimit = newAvailableLimit
}
