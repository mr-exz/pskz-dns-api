# pskz-dns-api

Go client for the [ps.kz](https://console.ps.kz) DNS GraphQL API.

## Quick start

```sh
git clone https://github.com/mr-exz/pskz-dns-api
cd pskz-dns-api

# List accounts linked to your user
PSKZ_EMAIL=you@example.com PSKZ_PASSWORD=yourpassword \
  go run ./cmd/list-accounts

# List DNS records for a specific account and zone
PSKZ_EMAIL=you@example.com PSKZ_PASSWORD=yourpassword \
PSKZ_ACCOUNT_ID=1001 PSKZ_ZONE=example.kz \
  go run ./cmd/list-records
```

## Installation

```sh
go get github.com/mr-exz/pskz-dns-api
```

## Authentication

Three options, pick one:

**Email + password** — recommended for automation:
```go
client, err := pskzdns.Login(ctx, "you@example.com", "password")
```

**Email + password + TOTP** — for accounts with two-factor authentication:
```go
client, err := pskzdns.LoginTwoFactor(ctx, "you@example.com", "password", "123456")
```

**Session cookie** — grab `connect.sid` from browser DevTools (Application → Cookies) after logging in to [console.ps.kz](https://console.ps.kz):
```go
client := pskzdns.New("s%3A...")
```

## Usage

```go
import pskzdns "github.com/mr-exz/pskz-dns-api"
```

A ps.kz user can have multiple billing accounts. You must select one before making DNS calls.

### List accounts

```go
accounts, err := client.ListAccounts(ctx)
if err != nil {
    log.Fatal(err)
}
for _, a := range accounts {
    fmt.Printf("%d  %s\n", a.ID, a.CompanyName)
}
```

### Select account

```go
if err := client.SelectAccount(ctx, 1001); err != nil {
    log.Fatal(err)
}
```

### List records

```go
records, err := client.ListRecords(ctx, "example.kz")
if err != nil {
    log.Fatal(err)
}
for _, r := range records {
    fmt.Printf("%s  %s  %d  %s\n", r.Name, r.Type, r.TTL, r.Value)
}
```

### Create a record

```go
record, err := client.CreateRecord(ctx, "example.kz", pskzdns.CreateRecordInput{
    Name:  "sub.example.kz",
    Type:  pskzdns.RecordTypeA,
    Value: "1.2.3.4",
    TTL:   180,
})
```

### Delete a record

```go
deleted, err := client.DeleteRecord(ctx, "example.kz", record.ID)
```

## Types

| Constant | Value |
|---|---|
| `RecordTypeA` | `A` |
| `RecordTypeAAAA` | `AAAA` |
| `RecordTypeCNAME` | `CNAME` |
| `RecordTypeMX` | `MX` |
| `RecordTypeTXT` | `TXT` |
| `RecordTypeNS` | `NS` |
| `RecordTypeSOA` | `SOA` |
| `RecordTypeSRV` | `SRV` |
| `RecordTypeCAA` | `CAA` |
