// Package apiclients provides a high-level client for IBM Verify API client management
// (Dynamic Client Registration). It wraps the Fern-generated apiclients client,
// handling token acquisition automatically so callers never manage access tokens directly.
//
// Note: The Fern-generated raw client always returns Body: nil for these endpoints
// because IBM's OpenAPI spec defines no response schemas for DCR operations. List,
// Get, and Create therefore use raw HTTP directly — the same approach used by apps
// and users. Delete still uses the generated typed client (no response body needed).
package apiclients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/apiclients"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
)

// Client manages IBM Verify API clients (Dynamic Client Registration).
type Client struct {
	tenantURL  string
	getToken   func(ctx context.Context) (string, error)
	httpClient *http.Client
}

// New returns an APIClients Client.
// getToken is a function that returns a valid bearer token.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{
		tenantURL:  tenantURL,
		getToken:   getToken,
		httpClient: &http.Client{},
	}
}

// newGenerated builds the Fern-generated client (used only for Delete).
func (c *Client) newGenerated(ctx context.Context) (*apiclients.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("apiclients: get token: %w", err)
	}
	h := make(http.Header)
	h.Set("Accept", "application/json")
	return apiclients.NewClient(&core.RequestOptions{
		BaseURL:    c.tenantURL,
		APIKey:     token,
		HTTPHeader: h,
	}), nil
}

// rawDo performs an authenticated HTTP request and returns the response body.
// Used for List, Get, and Create where the Fern raw client returns Body: nil.
func (c *Client) rawDo(ctx context.Context, method, path string, body any) ([]byte, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("apiclients: get token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("apiclients: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.tenantURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("apiclients: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apiclients: send request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("apiclients: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("apiclients: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// List returns all API clients as a raw JSON map slice.
func (c *Client) List(ctx context.Context, _ *generated.GetAPIClientsRequest) ([]map[string]interface{}, error) {
	body, err := c.rawDo(ctx, http.MethodGet, "/v1.0/apiclients", nil)
	if err != nil {
		return nil, fmt.Errorf("apiclients: list: %w", err)
	}
	var list []map[string]interface{}
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("apiclients: list: unmarshal: %w", err)
	}
	return list, nil
}

// Get returns a single API client by client ID as a raw JSON map.
func (c *Client) Get(ctx context.Context, clientID string) (map[string]interface{}, error) {
	body, err := c.rawDo(ctx, http.MethodGet, "/v1.0/apiclients/"+clientID, nil)
	if err != nil {
		return nil, fmt.Errorf("apiclients: get %s: %w", clientID, err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("apiclients: get %s: unmarshal: %w", clientID, err)
	}
	return m, nil
}

// Create registers a new API client and returns it as a raw JSON map.
func (c *Client) Create(ctx context.Context, req *generated.APIClientConfigRequest) (map[string]interface{}, error) {
	body, err := c.rawDo(ctx, http.MethodPost, "/v1.0/apiclients", req)
	if err != nil {
		return nil, fmt.Errorf("apiclients: create: %w", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("apiclients: create: unmarshal: %w", err)
	}
	return m, nil
}

// Delete removes an API client by client ID.
func (c *Client) Delete(ctx context.Context, clientID string) error {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return err
	}
	if err := cl.DeleteAPIClient(ctx, &generated.DeleteAPIClientRequest{
		ClientID: clientID,
	}); err != nil {
		return fmt.Errorf("apiclients: delete %s: %w", clientID, err)
	}
	return nil
}
