package pskzdns

// RecordType represents a DNS record type.
type RecordType string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeTXT   RecordType = "TXT"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeCAA   RecordType = "CAA"
)

// Record represents a DNS record.
type Record struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Type  RecordType `json:"type"`
	Value string     `json:"content"`
	TTL   int        `json:"ttl"`
}

// Account represents a ps.kz billing account.
type Account struct {
	ID          int    `json:"id"`
	CompanyName string `json:"companyname"`
	IsCurrent   bool   `json:"isCurrent"`
}

// CreateRecordInput holds the fields required to create a DNS record.
type CreateRecordInput struct {
	Name  string
	Type  RecordType
	Value string
	TTL   int
}
