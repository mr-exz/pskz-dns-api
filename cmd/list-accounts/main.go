package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pskzdns "github.com/mr-exz/pskz-dns-api"
)

func main() {
	email := os.Getenv("PSKZ_EMAIL")
	password := os.Getenv("PSKZ_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("set PSKZ_EMAIL and PSKZ_PASSWORD")
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

	fmt.Printf("%-10s  %s\n", "ID", "Company")
	fmt.Printf("%-10s  %s\n", "----------", "-------")
	for _, a := range accounts {
		current := ""
		if a.IsCurrent {
			current = " *"
		}
		fmt.Printf("%-10d  %s%s\n", a.ID, a.CompanyName, current)
	}
}
