// Package users provides a high-level client for IBM Verify user management (SCIM v2).
// It wraps the Fern-generated usersmanagementversion20 client, handling token
// acquisition automatically so callers never need to manage access tokens directly.
//
// Note: IBM Verify's SCIM responses contain type mismatches vs the OpenAPI spec
// (e.g. pwdChangedTime returned as string vs int64 in the spec). List and Get
// use raw HTTP to avoid generated unmarshalling failures. Write operations use
// the generated typed clients. When IBM updates their spec, swap back in one place.
package users

import (
	"context"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
	"github.com/ChrisVerde02/ibmverify-go/generated/usersmanagementversion20"
)

// Client manages IBM Verify users via the SCIM v2 API.
type Client struct {
	tenantURL  string
	getToken   func(ctx context.Context) (string, error)
	httpClient *http.Client
}

// New returns a Users Client.
// getToken is a function that returns a valid bearer token.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{
		tenantURL:  tenantURL,
		getToken:   getToken,
		httpClient: &http.Client{},
	}
}

// newGenerated builds the Fern-generated client with auth headers set.
func (c *Client) newGenerated(ctx context.Context) (*usersmanagementversion20.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("users: get token: %w", err)
	}
	h := make(http.Header)
	h.Set("Accept", "application/scim+json")
	h.Set("Content-Type", "application/scim+json")
	return usersmanagementversion20.NewClient(&core.RequestOptions{
		BaseURL:    c.tenantURL,
		APIKey:     token,
		HTTPHeader: h,
	}), nil
}

// rawGet performs a plain authenticated GET and returns the raw response body.
func (c *Client) rawGet(ctx context.Context, path string) ([]byte, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("users: get token: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.tenantURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("users: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/scim+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("users: send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("users: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("users: HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// List returns users as raw JSON maps.
// Pass nil to return all users, or set req.Filter for a SCIM filter expression
// (e.g. `userName eq "john"`).
func (c *Client) List(ctx context.Context, req *generated.GetUsersRequest) ([]map[string]interface{}, error) {
	path := "/v2.0/Users"
	if req != nil && req.Filter != nil && *req.Filter != "" {
		path += "?filter=" + url.QueryEscape(*req.Filter)
	}
	body, err := c.rawGet(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("users: list: %w", err)
	}
	// SCIM list response: {"Resources":[...], "totalResults":N}
	var wrapper struct {
		Resources []map[string]interface{} `json:"Resources"`
	}
	if json.Unmarshal(body, &wrapper) == nil && wrapper.Resources != nil {
		return wrapper.Resources, nil
	}
	var list []map[string]interface{}
	_ = json.Unmarshal(body, &list)
	return list, nil
}

// Get returns a single user by ID as a raw JSON map.
func (c *Client) Get(ctx context.Context, id string) (map[string]interface{}, error) {
	body, err := c.rawGet(ctx, "/v2.0/Users/"+id)
	if err != nil {
		return nil, fmt.Errorf("users: get %s: %w", id, err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("users: get %s: unmarshal: %w", id, err)
	}
	return m, nil
}

// Create creates a new user via raw HTTP to avoid typed-unmarshal failures
// caused by IBM's spec mismatch (e.g. pwdChangedTime returned as string, not int64).
// Returns the full SCIM response as a raw map.
func (c *Client) Create(ctx context.Context, req *generated.CreateUserRequest) (map[string]interface{}, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("users: get token: %w", err)
	}

	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		return nil, fmt.Errorf("users: create: marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.tenantURL+"/v2.0/Users", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("users: create: request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/scim+json")
	httpReq.Header.Set("Accept", "application/scim+json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("users: create: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("users: create: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("users: create: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var m map[string]interface{}
	if err := json.Unmarshal(respBody, &m); err != nil {
		return nil, fmt.Errorf("users: create: unmarshal: %w", err)
	}
	return m, nil
}

// Delete removes a user by ID.
func (c *Client) Delete(ctx context.Context, id string) error {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return err
	}
	if err := cl.DeleteUser0(ctx, &generated.DeleteUser0Request{
		ID: id,
	}); err != nil {
		return fmt.Errorf("users: delete %s: %w", id, err)
	}
	return nil
}
