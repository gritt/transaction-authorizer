package main

import (
	"encoding/json"
	"fmt"

	"github.com/gritt/transaction-authorizer/internal/core/domain"
)

type Output struct {
	Account    domain.Account
	Violations []string
}

func (o Output) MarshalJSON() ([]byte, error) {
	type outputWithAccount struct {
		Account    domain.Account `json:"account"`
		Violations []string       `json:"violations"`
	}
	type outputWithEmptyAccount struct {
		Account    struct{} `json:"account"`
		Violations []string `json:"violations"`
	}

	emptyAccount := domain.Account{}
	if o.Account == emptyAccount {
		return json.Marshal(&outputWithEmptyAccount{Violations: o.Violations})
	}
	return json.Marshal(&outputWithAccount{Account: o.Account, Violations: o.Violations})
}

func parseOutput(account domain.Account, errs []error) string {
	violations := getViolationsWith(errs)
	output := Output{
		Account:    account,
		Violations: violations,
	}

	data, err := json.Marshal(output)
	if err != nil {
		fmt.Println("failed to marshal output: ", err)
		return ""
	}

	return string(data)
}

func getViolationsWith(errs []error) []string {
	violations := []string{}
	if len(errs) > 0 {
		for _, err := range errs {
			if err == nil {
				continue
			}

			violations = append(violations, err.Error())
		}
	}
	return violations
}
