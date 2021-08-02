package service

import "github.com/gritt/transaction-authorizer/internal/core/domain"

type (
	AccountRepository interface {
		SaveAccount(cardStatus bool, availableLimit int) (domain.Account, error)
		FindAccount() (domain.Account, error)
		UpdateAccountLimit(newAvailableLimit int)
	}

	AccountService struct {
		repository AccountRepository
	}
)

func NewAccountService(repository AccountRepository) AccountService {
	return AccountService{repository: repository}
}

func (s AccountService) CreateAccount(cardStatus bool, availableLimit int) (domain.Account, error) {
	if account, err := s.repository.SaveAccount(cardStatus, availableLimit); err == nil {
		return account, nil
	}
	return domain.Account{}, domain.ErrAccountAlreadyInitialized
}

func (s AccountService) GetAccount() (domain.Account, error) {
	if account, err := s.repository.FindAccount(); err == nil {
		return account, nil
	}
	return domain.Account{}, domain.ErrAccountNotInitialized
}

func (s AccountService) SetAccountLimit(limit int) domain.Account {
	if limit < 0 {
		limit = 0
	}
	s.repository.UpdateAccountLimit(limit)

	account, _ := s.repository.FindAccount()

	return account
}
