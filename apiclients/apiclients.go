// Package apiclients provides a high-level client for IBM Verify API client management
// (Dynamic Client Registration). It wraps the Fern-generated apiclients client,
// handling token acquisition automatically so callers never manage access tokens directly.
//
// Note: IBM Verify's API client endpoints do not return response bodies — List, Get,
// and Create return the raw JSON via the RawResponse. Use WithRawResponse for full
// access to the response payload.
package apiclients

import (
	"context"
	"encoding/json"
	"fmt"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/apiclients"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
)

// Client manages IBM Verify API clients (Dynamic Client Registration).
type Client struct {
	tenantURL string
	getToken  func(ctx context.Context) (string, error)
}

// New returns an APIClients Client.
// getToken is a function that returns a valid bearer token.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{tenantURL: tenantURL, getToken: getToken}
}

func (c *Client) newGenerated(ctx context.Context) (*apiclients.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("apiclients: get token: %w", err)
	}
	return apiclients.NewClient(&core.RequestOptions{
		BaseURL: c.tenantURL,
		APIKey:  token,
	}), nil
}

// List returns all API clients as a raw JSON map slice.
// IBM Verify's spec does not define a typed response schema for this endpoint.
func (c *Client) List(ctx context.Context, req *generated.GetAPIClientsRequest) ([]map[string]interface{}, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &generated.GetAPIClientsRequest{}
	}
	raw, err := cl.WithRawResponse.GetAPIClients(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("apiclients: list: %w", err)
	}
	return decodeList(raw.Body)
}

// Get returns a single API client by client ID as a raw JSON map.
func (c *Client) Get(ctx context.Context, clientID string) (map[string]interface{}, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	raw, err := cl.WithRawResponse.GetAPIClient(ctx, &generated.GetAPIClientRequest{
		ClientID: clientID,
	})
	if err != nil {
		return nil, fmt.Errorf("apiclients: get %s: %w", clientID, err)
	}
	return decodeMap(raw.Body)
}

// Create registers a new API client and returns it as a raw JSON map.
func (c *Client) Create(ctx context.Context, req *generated.APIClientConfigRequest) (map[string]interface{}, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	raw, err := cl.WithRawResponse.CreateAPIClient(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("apiclients: create: %w", err)
	}
	return decodeMap(raw.Body)
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

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func decodeMap(v any) (map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	return m, json.Unmarshal(b, &m)
}

func decodeList(v any) ([]map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var list []map[string]interface{}
	return list, json.Unmarshal(b, &list)
}
