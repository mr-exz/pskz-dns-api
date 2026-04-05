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
	accountIDStr := os.Getenv("PSKZ_ACCOUNT_ID")
	zone := os.Getenv("PSKZ_ZONE")

	if email == "" || password == "" || accountIDStr == "" || zone == "" {
		log.Fatal("set PSKZ_EMAIL, PSKZ_PASSWORD, PSKZ_ACCOUNT_ID, PSKZ_ZONE")
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Fatal("invalid PSKZ_ACCOUNT_ID:", err)
	}

	ctx := context.Background()

	client, err := pskzdns.Login(ctx, email, password)
	if err != nil {
		log.Fatal("login:", err)
	}

	if err := client.SelectAccount(ctx, accountID); err != nil {
		log.Fatal("select account:", err)
	}

	records, err := client.ListRecords(ctx, zone)
	if err != nil {
		log.Fatal("list records:", err)
	}

	fmt.Printf("%-40s  %-6s  %-6s  %s\n", "NAME", "TYPE", "TTL", "VALUE")
	fmt.Printf("%-40s  %-6s  %-6s  %s\n", "----", "----", "---", "-----")
	for _, r := range records {
		fmt.Printf("%-40s  %-6s  %-6d  %s\n", r.Name, r.Type, r.TTL, r.Value)
	}
}
