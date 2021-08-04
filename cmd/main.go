package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/unknown/authorizer/internal/core/service"
	"github.com/unknown/authorizer/internal/repository"
)

func main() {
	memoryRepository := repository.NewMemoryRepository()
	accountService := service.NewAccountService(&memoryRepository)
	transactionService := service.NewTransactionService(&memoryRepository, accountService)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := parseInput(scanner.Text())

		if input.isCreateAccount() {
			account, err := accountService.CreateAccount(input.Account)
			output := parseOutput(account, []error{err})
			fmt.Println(output)
			continue
		}

		account, errs := transactionService.AuthorizeTransaction(input.Transaction)
		output := parseOutput(account, errs)
		fmt.Println(output)
	}
}
