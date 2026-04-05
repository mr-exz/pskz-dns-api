package pskzdns

import "context"

// ListRecords returns all DNS records for the given zone.
func (c *Client) ListRecords(ctx context.Context, zoneName string) ([]Record, error) {
	const query = `
		query ListRecords($zoneName: String!) {
			dns {
				zone(zoneName: $zoneName) {
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
// The mutation returns the updated Zone type, so the returned Record is
// reconstructed from the input (ID will be empty).
func (c *Client) CreateRecord(ctx context.Context, zoneName string, input CreateRecordInput) (*Record, error) {
	const query = `
		mutation CreateRecord($zoneName: String!, $name: String!, $type: RecordType!, $value: String!, $ttl: Int!) {
			dns {
				record {
					create(zoneName: $zoneName, createData: {name: $name, type: $type, value: $value, ttl: $ttl}) {
						name
					}
				}
			}
		}`

	type response struct {
		DNS struct {
			Record struct {
				Create struct {
					Name string `json:"name"`
				} `json:"create"`
			} `json:"record"`
		} `json:"dns"`
	}

	if _, err := do[response](ctx, c, gqlRequest{
		Query: query,
		Variables: map[string]any{
			"zoneName": zoneName,
			"name":     input.Name,
			"type":     input.Type,
			"value":    input.Value,
			"ttl":      input.TTL,
		},
	}); err != nil {
		return nil, err
	}

	return &Record{
		Name:  input.Name,
		Type:  input.Type,
		Value: input.Value,
		TTL:   input.TTL,
	}, nil
}

// DeleteRecord deletes the DNS record with the given ID from the zone.
func (c *Client) DeleteRecord(ctx context.Context, zoneName, recordID string) error {
	const query = `
		mutation DeleteRecord($zoneName: String!, $recordId: String!) {
			dns {
				record {
					delete(zoneName: $zoneName, recordId: $recordId) {
						name
					}
				}
			}
		}`

	type response struct {
		DNS struct {
			Record struct {
				Delete struct {
					Name string `json:"name"`
				} `json:"delete"`
			} `json:"record"`
		} `json:"dns"`
	}

	_, err := do[response](ctx, c, gqlRequest{
		Query: query,
		Variables: map[string]any{
			"zoneName": zoneName,
			"recordId": recordID,
		},
	})
	return err
}
