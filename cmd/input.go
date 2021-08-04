package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/unknown/authorizer/internal/core/domain"
)

type Input struct {
	Transaction domain.Transaction `json:"transaction"`
	Account     domain.Account     `json:"account"`
}

func (o Input) isCreateAccount() bool {
	return o.Account != domain.Account{}
}

func parseInput(JSON string) Input {
	operation := Input{}
	if err := json.Unmarshal([]byte(JSON), &operation); err != nil {
		fmt.Println("failed to parse operation", err)
		os.Exit(1)
	}
	return operation
}
