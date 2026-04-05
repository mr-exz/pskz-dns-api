package pskzdns

import "context"

// ListRecords returns all DNS records for the given zone.
func (c *Client) ListRecords(ctx context.Context, zoneName string) ([]Record, error) {
	const query = `
		query ListRecords($zoneName: String!) {
			dns {
				zone(name: $zoneName) {
					records {
						id
						name
						type
						value
						ttl
					}
				}
			}
		}`

	type response struct {
		DNS struct {
			Zone struct {
				Records []Record `json:"records"`
			} `json:"zone"`
		} `json:"dns"`
	}

	data, err := do[response](ctx, c, gqlRequest{
		Query:     query,
		Variables: map[string]any{"zoneName": zoneName},
	})
	if err != nil {
		return nil, err
	}

	return data.DNS.Zone.Records, nil
}

// CreateRecord creates a new DNS record in the given zone.
func (c *Client) CreateRecord(ctx context.Context, zoneName string, input CreateRecordInput) (*Record, error) {
	const query = `
		mutation CreateRecord($zoneName: String!, $name: String!, $type: RecordType!, $value: String!, $ttl: Int!) {
			dns {
				record {
					create(zoneName: $zoneName, createData: {name: $name, type: $type, value: $value, ttl: $ttl}) {
						id
						name
						type
						value
						ttl
					}
				}
			}
		}`

	type response struct {
		DNS struct {
			Record struct {
				Create Record `json:"create"`
			} `json:"record"`
		} `json:"dns"`
	}

	data, err := do[response](ctx, c, gqlRequest{
		Query: query,
		Variables: map[string]any{
			"zoneName": zoneName,
			"name":     input.Name,
			"type":     input.Type,
			"value":    input.Value,
			"ttl":      input.TTL,
		},
	})
	if err != nil {
		return nil, err
	}

	r := data.DNS.Record.Create
	return &r, nil
}

// DeleteRecord deletes the DNS record with the given ID from the zone.
func (c *Client) DeleteRecord(ctx context.Context, zoneName, recordID string) (*Record, error) {
	const query = `
		mutation DeleteRecord($zoneName: String!, $recordId: String!) {
			dns {
				record {
					delete(zoneName: $zoneName, recordId: $recordId) {
						id
						name
						type
						value
						ttl
					}
				}
			}
		}`

	type response struct {
		DNS struct {
			Record struct {
				Delete Record `json:"delete"`
			} `json:"record"`
		} `json:"dns"`
	}

	data, err := do[response](ctx, c, gqlRequest{
		Query: query,
		Variables: map[string]any{
			"zoneName": zoneName,
			"recordId": recordID,
		},
	})
	if err != nil {
		return nil, err
	}

	r := data.DNS.Record.Delete
	return &r, nil
}
