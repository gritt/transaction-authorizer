package service

import (
	"time"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
)

type (
	AccountServicer interface {
		GetAccount() (domain.Account, error)
		SetAccountLimit(newAvailableLimit int) domain.Account
	}

	TransactionRepository interface {
		SaveTransaction(domain.Transaction)
		FindTransactionsAfter(time.Time) []domain.Transaction
	}

	TransactionService struct {
		repository     TransactionRepository
		accountService AccountServicer
	}
)

func NewTransactionService(repository TransactionRepository, accountManager AccountServicer) TransactionService {
	return TransactionService{
		repository:     repository,
		accountService: accountManager,
	}
}

func (s TransactionService) AuthorizeTransaction(transaction domain.Transaction) (domain.Account, []error) {
	account, err := s.accountService.GetAccount()
	if err != nil {
		return domain.Account{}, []error{err}
	}

	if !account.ActiveCard {
		return account, []error{domain.ErrCardNotActive}
	}

	errors := []error{}
	if account.AvailableLimit < transaction.Amount {
		errors = append(errors, domain.ErrInsufficientLimit)
	}

	twoMinutesAgo := transaction.CreatedAt.UTC().Add(-2 * time.Minute)
	pastTransactions := s.repository.FindTransactionsAfter(twoMinutesAgo)
	if len(pastTransactions) >= 3 {
		errors = append(errors, domain.ErrHighFrequencySmallInterval)
	}

	for _, pastTransaction := range pastTransactions {
		if pastTransaction.Amount == transaction.Amount && pastTransaction.Merchant == transaction.Merchant {
			errors = append(errors, domain.ErrDoubleTransaction)
			break
		}
	}

	if len(errors) >= 1 {
		return account, errors
	}

	s.repository.SaveTransaction(transaction)

	limit := account.AvailableLimit - transaction.Amount
	updatedAccount := s.accountService.SetAccountLimit(limit)

	return updatedAccount, []error{}
}
