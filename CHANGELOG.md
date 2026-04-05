# Changelog

## [0.0.1]

### Added
- `New(sessionID)` — create a client from a `connect.sid` browser session cookie
- `Login(ctx, email, password)` — authenticate programmatically with email and password
- `LoginTwoFactor(ctx, email, password, totpCode)` — authenticate with email, password, and TOTP code for 2FA-enabled accounts
- `ListAccounts(ctx)` — list all billing accounts linked to the authenticated user
- `SelectAccount(ctx, id)` — switch the client to operate under a specific account ID; required before DNS calls
- `ListRecords(ctx, zoneName)` — list all DNS records in a zone
- `CreateRecord(ctx, zoneName, input)` — create a DNS record in a zone
- `DeleteRecord(ctx, zoneName, recordID)` — delete a DNS record by ID
- `RecordType` constants: A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, CAA
