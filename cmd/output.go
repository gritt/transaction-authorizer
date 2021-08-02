package main

// TODO map to:
//  {"account":{"active-card":true,"available-limit":100},"violations":[]}
type operationLog struct {
	Account struct {
		ActiveCard     bool `json:"active-card"`
		AvailableLimit int  `json:"available-limit"`
	} `json:"account"`
	Violations []string `json:"violations"`
}
