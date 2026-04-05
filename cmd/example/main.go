package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	pskzdns "github.com/mr-exz/pskz-dns-api"
)

func main() {
	email := os.Getenv("PSKZ_EMAIL")
	password := os.Getenv("PSKZ_PASSWORD")
	zone := os.Getenv("PSKZ_ZONE")
	accountIDStr := os.Getenv("PSKZ_ACCOUNT_ID")

	if email == "" || password == "" || zone == "" {
		log.Fatal("set PSKZ_EMAIL, PSKZ_PASSWORD, PSKZ_ZONE")
	}

	ctx := context.Background()

	client, err := pskzdns.Login(ctx, email, password)
	if err != nil {
		log.Fatal("login:", err)
	}

	accounts, err := client.ListAccounts(ctx)
	if err != nil {
		log.Fatal("list accounts:", err)
	}

	fmt.Println("accounts:")
	for _, a := range accounts {
		current := ""
		if a.IsCurrent {
			current = " (current)"
		}
		fmt.Printf("  %d  %s%s\n", a.ID, a.CompanyName, current)
	}

	// Select account: use PSKZ_ACCOUNT_ID if set, otherwise the first account.
	accountID := accounts[0].ID
	if accountIDStr != "" {
		accountID, err = strconv.Atoi(accountIDStr)
		if err != nil {
			log.Fatal("invalid PSKZ_ACCOUNT_ID:", err)
		}
	}

	fmt.Printf("\nusing account %d\n", accountID)
	if err := client.SelectAccount(ctx, accountID); err != nil {
		log.Fatal("select account:", err)
	}

	records, err := client.ListRecords(ctx, zone)
	if err != nil {
		log.Fatal("list records:", err)
	}

	fmt.Printf("\nrecords in %s (%d):\n", zone, len(records))
	for _, r := range records {
		fmt.Printf("  %-40s %-6s %-6d %s\n", r.Name, r.Type, r.TTL, r.Value)
	}
}
