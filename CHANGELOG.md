# Changelog

## [0.0.5]

<!-- Prepare for next release: remove this line and write your release notes -->


## [0.0.4]

### Fixed
- `ListRecords`: changed query path from `dns.zone(name).records` to `dns.record.list(zoneName)` to match the mutation path convention used by create/delete


## [0.0.3]

### Fixed
- `ListRecords`: reverted field name back to `value` — `content` was incorrect; `Record` type uses `value`, only `CreateRecord`/`DeleteRecord` mutation responses were affected by the Zone type mismatch
- `Record.Value`: reverted JSON tag back to `json:"value"` to match the ps.kz API schema


## [0.0.2]

### Fixed
- `ListRecords`: changed requested field from `value` to `content` to match ps.kz API schema; records were returned with blank values previously
- `CreateRecord`: mutation returns `Zone` type, not `Record`; response parsing updated to request only `name` from the zone; returned `Record` is now reconstructed from the input (ID will be empty)
- `DeleteRecord`: same `Zone` return type fix; signature simplified to return `error` only
- `doAt`: return an error instead of a nil pointer when the API response contains no `data` field


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
