package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unknown/authorizer/internal/core/domain"
)

func TestCreateAccount(t *testing.T) {
	givenErr := errors.New("repository error")
	givenAccount := domain.Account{
		ActiveCard:     false,
		AvailableLimit: 100,
	}

	testCases := map[string]func(*testing.T, *accountRepositoryMock){
		"should create account with success": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			accountRepositoryMock.On("SaveAccount", givenAccount).Return(givenAccount, nil)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account, err := accountService.CreateAccount(givenAccount)

			// 	then
			assert.Equal(t, givenAccount, account)
			assert.NoError(t, err)
		},
		"should return error when repository fails to save": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			accountRepositoryMock.On("SaveAccount", givenAccount).Return(domain.Account{}, givenErr)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account, err := accountService.CreateAccount(givenAccount)

			// 	then
			assert.Empty(t, account)
			assert.EqualError(t, err, domain.ErrAccountAlreadyInitialized.Error())
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			accountRepositoryMock := new(accountRepositoryMock)

			run(t, accountRepositoryMock)

			accountRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestGetAccount(t *testing.T) {
	givenErr := errors.New("repository error")
	givenAccount := domain.Account{
		ActiveCard:     false,
		AvailableLimit: 100,
	}

	testCases := map[string]func(*testing.T, *accountRepositoryMock){
		"should get account with success": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			accountRepositoryMock.On("FindAccount").Return(givenAccount, nil)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account, err := accountService.GetAccount()

			// 	then
			assert.Equal(t, givenAccount, account)
			assert.NoError(t, err)
		},
		"should return error when repository fails to find": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			accountRepositoryMock.On("FindAccount").Return(domain.Account{}, givenErr)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account, err := accountService.GetAccount()

			// 	then
			assert.Empty(t, account)
			assert.EqualError(t, err, domain.ErrAccountNotInitialized.Error())
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			accountRepositoryMock := new(accountRepositoryMock)

			run(t, accountRepositoryMock)

			accountRepositoryMock.AssertExpectations(t)
		})
	}
}

func TestSetAccountLimit(t *testing.T) {
	testCases := map[string]func(*testing.T, *accountRepositoryMock){
		"should update limit and return account": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			givenAccount := domain.Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			accountRepositoryMock.On("UpdateAccountLimit", 100)
			accountRepositoryMock.On("FindAccount").Return(givenAccount, nil)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account := accountService.SetAccountLimit(100)

			// 	then
			assert.Equal(t, givenAccount, account)
		},
		"should update limit with zero when negative value and return account": func(t *testing.T, accountRepositoryMock *accountRepositoryMock) {
			// 	given
			givenAccount := domain.Account{
				ActiveCard:     true,
				AvailableLimit: 0,
			}

			accountRepositoryMock.On("UpdateAccountLimit", 0)
			accountRepositoryMock.On("FindAccount").Return(givenAccount, nil)

			accountService := NewAccountService(accountRepositoryMock)

			// 	when
			account := accountService.SetAccountLimit(-1)

			// 	then
			assert.Equal(t, givenAccount, account)
		},
	}

	for name, run := range testCases {
		t.Run(name, func(t *testing.T) {
			accountRepositoryMock := new(accountRepositoryMock)

			run(t, accountRepositoryMock)

			accountRepositoryMock.AssertExpectations(t)
		})
	}
}
