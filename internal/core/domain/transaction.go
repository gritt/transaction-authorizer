package domain

import "time"

type Transaction struct {
	Amount    int
	Merchant  string
	CreatedAt time.Time
}
