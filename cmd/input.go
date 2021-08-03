package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
)

type Operation struct {
	Transaction domain.Transaction `json:"transaction"`
	Account     domain.Account     `json:"account"`
}

func (o Operation) isCreateAccount() bool {
	return o.Account != domain.Account{}
}

func parseOperation(JSON string) Operation {
	operation := Operation{}
	if err := json.Unmarshal([]byte(JSON), &operation); err != nil {
		fmt.Println("failed to parse operation", err)
		os.Exit(1)
	}
	return operation
}
