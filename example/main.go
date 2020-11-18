// The simple example demonstrates fetching a list of accounts
package main

import (
	"context"
	"fmt"

	"form3.tech/go-form3/form3"
)

func fetchAccounts() (*form3.AccountDetailsListResponse, error) {
	client := form3.NewClient(nil)
	accounts, _, err := client.Accounts.List(context.Background(), &form3.ListOptions{PageNumber: 1, PageSize: 50})
	return accounts, err
}

func main() {
	accounts, err := fetchAccounts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, account := range accounts.Data {
		fmt.Printf("%v. %v\n", i+1, account.ID)
	}
}
