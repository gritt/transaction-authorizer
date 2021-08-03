package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gritt/transaction-authorizer/internal/core/service"
	"github.com/gritt/transaction-authorizer/internal/repository"
)

func main() {
	memoryRepository := repository.NewMemoryRepository()
	accountService := service.NewAccountService(&memoryRepository)
	transactionService := service.NewTransactionService(&memoryRepository, accountService)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		operation := parseOperation(scanner.Text())

		if operation.isCreateAccount() {
			account, err := accountService.CreateAccount(operation.Account)
			txt := fmt.Sprintf("%v : %s", account, err)
			fmt.Println(txt)
			continue
		}

		transaction, err := transactionService.AuthorizeTransaction(operation.Transaction)
		txt := fmt.Sprintf("%v : %s", transaction, err)
		fmt.Println(txt)
	}
}
