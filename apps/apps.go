// Package apps provides a high-level client for IBM Verify application management.
// It wraps the Fern-generated applicationaccess client, handling token acquisition
// automatically so callers never need to manage access tokens directly.
package apps

import (
	"context"
	"fmt"
	"net/http"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/applicationaccess"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
)

// Client manages IBM Verify applications.
type Client struct {
	tenantURL string
	getToken  func(ctx context.Context) (string, error)
}

// New returns an Apps Client.
// getToken is a function that returns a valid bearer token — typically
// wrapping client.Token.ClientCredentials(ctx).AccessToken.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{tenantURL: tenantURL, getToken: getToken}
}

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

// List searches applications. Pass nil for defaults (no filter, no pagination).
func (c *Client) List(ctx context.Context, req *generated.SearchApplicationsRequest) (*generated.SearchAdminApplicationWithoutProvResponseBean, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &generated.SearchApplicationsRequest{}
	}
	result, err := cl.SearchApplications(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("apps: list: %w", err)
	}
	return result, nil
}

// Get returns a single application by ID.
func (c *Client) Get(ctx context.Context, applicationID string) (*generated.ApplicationDetailsResponseBean, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	result, err := cl.GetApplication(ctx, &generated.GetApplicationRequest{
		ApplicationID: applicationID,
	})
	if err != nil {
		return nil, fmt.Errorf("apps: get %s: %w", applicationID, err)
	}
	return result, nil
}

// Create registers a new application.
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
