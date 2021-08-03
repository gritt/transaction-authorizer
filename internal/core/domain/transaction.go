package domain

import "time"

type Transaction struct {
	Amount    int       `json:"amount"`
	Merchant  string    `json:"merchant"`
	CreatedAt time.Time `json:"time"`
}
