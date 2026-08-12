package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// TokenClient provides OAuth token operations against IBM Verify.
// Access it via Client.Token.
type TokenClient struct {
	c *Client
}

// ClientCredentialsResult holds the response from a client credentials grant.
type ClientCredentialsResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// ExchangeResult holds the response from a token exchange grant.
type ExchangeResult struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn       int64  `json:"expires_in"`
	GrantID         string `json:"grant_id"`
	IssuedTokenType string `json:"issued_token_type"`
	Scope           string `json:"scope"`
	TokenType       string `json:"token_type"`
}

// IntrospectResult holds token metadata returned by IBM Verify.
type IntrospectResult struct {
	Active            bool   `json:"active"`
	ClientID          string `json:"client_id"`
	Username          string `json:"username"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	Subject           string `json:"sub"`
	Scope             string `json:"scope"`
	TokenType         string `json:"token_type"`
	Issuer            string `json:"iss"`
	IssuedAt          int64  `json:"iat"`
	ExpiresAt         int64  `json:"exp"`
}

// ClientCredentials obtains an access token using the OAuth 2.0 client
// credentials grant configured on the parent Client.
//
//	POST /v1.0/endpoint/default/token
func (t *TokenClient) ClientCredentials(ctx context.Context) (*ClientCredentialsResult, error) {
	if t.c.clientID == "" {
		return nil, fmt.Errorf("client credentials: client ID not configured — use WithClientCredentials()")
	}
	if t.c.clientSecret == "" {
		return nil, fmt.Errorf("client credentials: client secret not configured — use WithClientCredentials()")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", t.c.clientID)
	form.Set("client_secret", t.c.clientSecret)

	body, err := t.c.postForm(ctx, "/v1.0/endpoint/default/token", form, "")
	if err != nil {
		return nil, fmt.Errorf("client credentials: %w", err)
	}

	var result ClientCredentialsResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("client credentials: decode response: %w", err)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("client credentials: response did not contain an access_token")
	}
	return &result, nil
}

// Exchange exchanges a subject token (e.g. a signed JWT) for an IBM Verify
// access token using the OAuth 2.0 Token Exchange grant (RFC 8693).
//
//	POST /oauth2/token
func (t *TokenClient) Exchange(
	ctx context.Context,
	subjectToken string,
	subjectTokenType string,
) (*ExchangeResult, error) {
	if strings.TrimSpace(subjectToken) == "" {
		return nil, fmt.Errorf("token exchange: subject token cannot be empty")
	}
	if strings.TrimSpace(subjectTokenType) == "" {
		return nil, fmt.Errorf("token exchange: subject token type cannot be empty")
	}
	if t.c.clientID == "" {
		return nil, fmt.Errorf("token exchange: client ID not configured — use WithClientCredentials()")
	}

	form := url.Values{}
	form.Set("client_id", t.c.clientID)
	form.Set("client_secret", t.c.clientSecret)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	form.Set("subject_token", subjectToken)
	form.Set("subject_token_type", subjectTokenType)

	body, err := t.c.postForm(ctx, "/oauth2/token", form, "")
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}

	var result ExchangeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("token exchange: decode response: %w", err)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("token exchange: response did not contain an access_token")
	}
	return &result, nil
}

// Introspect asks IBM Verify for metadata about an access token.
//
//	POST /oauth2/introspect
func (t *TokenClient) Introspect(ctx context.Context, token string) (*IntrospectResult, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("introspect: token cannot be empty")
	}
	if t.c.clientID == "" {
		return nil, fmt.Errorf("introspect: client ID not configured — use WithClientCredentials()")
	}

	form := url.Values{}
	form.Set("client_id", t.c.clientID)
	form.Set("client_secret", t.c.clientSecret)
	form.Set("token", token)
	form.Set("token_type_hint", "access_token")

	body, err := t.c.postForm(ctx, "/oauth2/introspect", form, "")
	if err != nil {
		return nil, fmt.Errorf("introspect: %w", err)
	}

	var result IntrospectResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("introspect: decode response: %w", err)
	}
	return &result, nil
}
