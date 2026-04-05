package pskzdns

import "context"

// ListAccounts returns all accounts linked to the authenticated user.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	const query = `{ account { list { id companyname isCurrent } } }`

	type response struct {
		Account struct {
			List []Account `json:"list"`
		} `json:"account"`
	}

	data, err := doAt[response](ctx, c, accountEndpoint, gqlRequest{Query: query})
	if err != nil {
		return nil, err
	}

	return data.Account.List, nil
}

// SelectAccount switches the client to operate under the given account ID.
// Must be called before making DNS calls when working with a specific account.
func (c *Client) SelectAccount(ctx context.Context, id int) error {
	type response struct {
		Account struct {
			SetAccount struct {
				ID int `json:"id"`
			} `json:"setAccount"`
		} `json:"account"`
	}

	_, err := doAt[response](ctx, c, accountEndpoint, gqlRequest{
		Query:     `mutation SelectAccount($id: Int!) { account { setAccount(id: $id) { id } } }`,
		Variables: map[string]any{"id": id},
	})
	return err
}
