package main

import (
	"os"
)

var defaultStdin = os.Stdin
var testFile *os.File

func Example_main_when_create_account() {
	setup("../test/create_account")
	defer teardown()

	main()

	// Output:
	// {"account":{"active-card":false,"available-limit":750},"violations":[]}
}

func Example_main_when_account_not_initialized() {
	setup("../test/account_not_initialized")
	defer teardown()

	main()

	// Output:
	// {"account":{},"violations":["account-not-initialized"]}
}

func Example_main_when_account_card_not_active() {
	setup("../test/card_not_active")
	defer teardown()

	main()

	// Output:
	// {"account":{"active-card":false,"available-limit":100},"violations":[]}
	// {"account":{"active-card":false,"available-limit":100},"violations":["card-not-active"]}
}

func Example_main_when_has_multiple_violations() {
	setup("../test/multiple_violations")
	defer teardown()

	main()

	// Output:
	// {"account":{"active-card":true,"available-limit":100},"violations":[]}
	// {"account":{"active-card":true,"available-limit":90},"violations":[]}
	// {"account":{"active-card":true,"available-limit":70},"violations":[]}
	// {"account":{"active-card":true,"available-limit":65},"violations":[]}
	// {"account":{"active-card":true,"available-limit":65},"violations":["high-frequency-small-interval","double-transaction"]}
	// {"account":{"active-card":true,"available-limit":65},"violations":["insufficient-limit","high-frequency-small-interval"]}
	// {"account":{"active-card":true,"available-limit":65},"violations":["insufficient-limit","high-frequency-small-interval"]}
	// {"account":{"active-card":true,"available-limit":50},"violations":[]}
}

func setup(path string) {
	testFile, _ = os.Open(path)
	os.Stdin = testFile
}

func teardown() {
	os.Stdin = defaultStdin
	testFile.Close()
}
