package pskzdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

const (
	defaultEndpoint    = "https://console.ps.kz/dns/graphql"
	authEndpoint       = "https://console.ps.kz/auth/graphql"
	accountEndpoint    = "https://console.ps.kz/account/graphql"
)

// Client is a ps.kz DNS API client.
type Client struct {
	http     *http.Client
	endpoint string
}

// New creates a Client using a connect.sid session cookie value.
func New(sessionID string) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	u, _ := url.Parse(defaultEndpoint)
	jar.SetCookies(u, []*http.Cookie{{Name: "connect.sid", Value: sessionID}})
	return &Client{
		http:     &http.Client{Jar: jar},
		endpoint: defaultEndpoint,
	}
}

// Login authenticates with email and password and returns an authenticated Client.
// For accounts with two-factor authentication enabled, use LoginTwoFactor instead.
func Login(ctx context.Context, email, password string) (*Client, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c := &Client{http: &http.Client{Jar: jar}, endpoint: defaultEndpoint}

	authState, err := guestLogin(ctx, c, email, password)
	if err != nil {
		return nil, err
	}
	if authState == "twoFactor" {
		return nil, ErrTwoFactorRequired
	}

	return c, nil
}

// LoginTwoFactor authenticates with email, password, and a TOTP code for accounts
// with two-factor authentication enabled.
func LoginTwoFactor(ctx context.Context, email, password, totpCode string) (*Client, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c := &Client{http: &http.Client{Jar: jar}, endpoint: defaultEndpoint}

	authState, err := guestLogin(ctx, c, email, password)
	if err != nil {
		return nil, err
	}
	if authState != "twoFactor" {
		// 2FA not required — already authenticated.
		return c, nil
	}

	const query = `
		mutation TwoFactor($token: String!) {
			auth {
				guestLoginTwoFactor(token: $token, rememberDevice: false) {
					authState
				}
			}
		}`

	type response struct {
		Auth struct {
			GuestLoginTwoFactor struct {
				AuthState string `json:"authState"`
			} `json:"guestLoginTwoFactor"`
		} `json:"auth"`
	}

	if _, err := doAt[response](ctx, c, authEndpoint, gqlRequest{
		Query:     query,
		Variables: map[string]any{"token": totpCode},
	}); err != nil {
		return nil, err
	}

	return c, nil
}

// ErrTwoFactorRequired is returned by Login when the account has 2FA enabled.
var ErrTwoFactorRequired = fmt.Errorf("pskzdns: two-factor authentication required")

func guestLogin(ctx context.Context, c *Client, email, password string) (authState string, err error) {
	const query = `
		mutation Login($email: EmailAddress!, $password: String!) {
			auth {
				guestLogin(email: $email, password: $password, remember: true) {
					authState
				}
			}
		}`

	type response struct {
		Auth struct {
			GuestLogin struct {
				AuthState string `json:"authState"`
			} `json:"guestLogin"`
		} `json:"auth"`
	}

	data, err := doAt[response](ctx, c, authEndpoint, gqlRequest{
		Query:     query,
		Variables: map[string]any{"email": email, "password": password},
	})
	if err != nil {
		return "", err
	}

	return data.Auth.GuestLogin.AuthState, nil
}


type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type gqlResponse[T any] struct {
	Data   *T         `json:"data"`
	Errors []gqlError `json:"errors"`
}

type gqlError struct {
	Message string `json:"message"`
}

func do[T any](ctx context.Context, c *Client, req gqlRequest) (*T, error) {
	return doAt[T](ctx, c, c.endpoint, req)
}

func doAt[T any](ctx context.Context, c *Client, endpoint string, req gqlRequest) (*T, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gqlResp gqlResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, err
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("graphql: %s", gqlResp.Errors[0].Message)
	}

	return gqlResp.Data, nil
}
