// Package apps provides a high-level client for IBM Verify application management.
// It wraps the Fern-generated applicationaccess client, handling token acquisition
// automatically so callers never need to manage access tokens directly.
//
// Note: IBM Verify's application API responses contain type mismatches vs the OpenAPI
// spec (e.g. applicationState returned as bool vs string in the spec). List and Get
// use raw HTTP to avoid generated unmarshalling failures. Write operations (Create,
// Delete) use the generated typed clients since we control the request payload.
// When IBM updates their spec, swap List/Get back to typed calls in one place here.
package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/applicationaccess"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
)

// Client manages IBM Verify applications.
type Client struct {
	tenantURL  string
	getToken   func(ctx context.Context) (string, error)
	httpClient *http.Client
}

// New returns an Apps Client.
// getToken is a function that returns a valid bearer token.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{
		tenantURL:  tenantURL,
		getToken:   getToken,
		httpClient: &http.Client{},
	}
}

// newGenerated builds the Fern-generated client with auth headers set.
func (c *Client) newGenerated(ctx context.Context) (*applicationaccess.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("apps: get token: %w", err)
	}
	h := make(http.Header)
	h.Set("Accept", "application/json")
	return applicationaccess.NewClient(&core.RequestOptions{
		BaseURL:    c.tenantURL,
		APIKey:     token,
		HTTPHeader: h,
	}), nil
}

// rawGet performs a plain authenticated GET and returns the raw response body.
// Used for endpoints where IBM's live response doesn't match their OpenAPI spec types.
func (c *Client) rawGet(ctx context.Context, path string) ([]byte, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("apps: get token: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.tenantURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("apps: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apps: send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("apps: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("apps: HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// List returns all applications as raw JSON maps.
// IBM returns {"_embedded":{"applications":[...]}} — this unwraps that automatically.
func (c *Client) List(ctx context.Context, _ *generated.SearchApplicationsRequest) ([]map[string]interface{}, error) {
	body, err := c.rawGet(ctx, "/v1.0/applications")
	if err != nil {
		return nil, fmt.Errorf("apps: list: %w", err)
	}
	var wrapper struct {
		Embedded *struct {
			Applications []map[string]interface{} `json:"applications"`
		} `json:"_embedded"`
	}
	if json.Unmarshal(body, &wrapper) == nil && wrapper.Embedded != nil {
		return wrapper.Embedded.Applications, nil
	}
	var list []map[string]interface{}
	_ = json.Unmarshal(body, &list)
	return list, nil
}

// Get returns a single application by ID as a raw JSON map.
func (c *Client) Get(ctx context.Context, applicationID string) (map[string]interface{}, error) {
	body, err := c.rawGet(ctx, "/v1.0/applications/"+applicationID)
	if err != nil {
		return nil, fmt.Errorf("apps: get %s: %w", applicationID, err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("apps: get %s: unmarshal: %w", applicationID, err)
	}
	return m, nil
}

// Create registers a new application using the generated typed client.
func (c *Client) Create(ctx context.Context, req *generated.ApplicationRequestBean) (*generated.PostApplicationResponseBean, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	result, err := cl.CreateApplication(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("apps: create: %w", err)
	}
	return result, nil
}

// Delete removes an application by ID.
func (c *Client) Delete(ctx context.Context, applicationID string) error {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return err
	}
	if err := cl.DeleteApplication(ctx, &generated.DeleteApplicationRequest{
		ApplicationID: applicationID,
	}); err != nil {
		return fmt.Errorf("apps: delete %s: %w", applicationID, err)
	}
	return nil
}
