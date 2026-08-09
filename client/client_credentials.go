package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ClientCredentialsRequest contains the values sent to IBM Verify.
type ClientCredentialsRequest struct {
	TenantURL    string
	ClientID     string
	ClientSecret string
}

// ClientCredentialsResponse represents the token returned by IBM Verify.
type ClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetClientCredentialsToken fetches an access token from IBM Verify using the
// OAuth 2.0 client credentials grant.
//
//	POST /v1.0/endpoint/default/token
func GetClientCredentialsToken(
	ctx context.Context,
	request ClientCredentialsRequest,
) (*ClientCredentialsResponse, error) {
	if strings.TrimSpace(request.TenantURL) == "" {
		return nil, errors.New("tenant URL cannot be empty")
	}

	if strings.TrimSpace(request.ClientID) == "" {
		return nil, errors.New("client ID cannot be empty")
	}

	if strings.TrimSpace(request.ClientSecret) == "" {
		return nil, errors.New("client secret cannot be empty")
	}

	endpoint := strings.TrimRight(request.TenantURL, "/") + "/v1.0/endpoint/default/token"

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", request.ClientID)
	form.Set("client_secret", request.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create client credentials request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send client credentials request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read client credentials response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"IBM Verify client credentials failed with HTTP %d: %s",
			resp.StatusCode, string(body),
		)
	}

	var result ClientCredentialsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode client credentials response: %w", err)
	}

	if result.AccessToken == "" {
		return nil, errors.New("IBM Verify response did not contain an access_token")
	}

	return &result, nil
}
