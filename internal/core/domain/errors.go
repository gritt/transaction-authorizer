package domain

import "github.com/pkg/errors"

var (
	ErrAccountNotInitialized      = errors.New("account-not-initialized")
	ErrAccountAlreadyInitialized  = errors.New("account-already-initialized")
	ErrCardNotActive              = errors.New("card-not-active")
	ErrInsufficientLimit          = errors.New("insufficient-limit")
	ErrHighFrequencySmallInterval = errors.New("high-frequency-small-interval")
	ErrDoubleTransaction          = errors.New("double-transaction")
)
