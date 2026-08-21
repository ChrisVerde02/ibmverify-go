// Package users provides a high-level client for IBM Verify user management (SCIM v2).
// It wraps the Fern-generated usersmanagementversion20 client, handling token
// acquisition automatically so callers never need to manage access tokens directly.
package users

import (
	"context"
	"fmt"

	generated "github.com/ChrisVerde02/ibmverify-go/generated"
	"github.com/ChrisVerde02/ibmverify-go/generated/core"
	"github.com/ChrisVerde02/ibmverify-go/generated/usersmanagementversion20"
)

// Client manages IBM Verify users via the SCIM v2 API.
type Client struct {
	tenantURL string
	getToken  func(ctx context.Context) (string, error)
}

// New returns a Users Client.
// getToken is a function that returns a valid bearer token.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
	return &Client{tenantURL: tenantURL, getToken: getToken}
}

func (c *Client) newGenerated(ctx context.Context) (*usersmanagementversion20.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("users: get token: %w", err)
	}
	return usersmanagementversion20.NewClient(&core.RequestOptions{
		BaseURL: c.tenantURL,
		APIKey:  token,
	}), nil
}

// List returns users matching the request filter. Pass nil for all users.
func (c *Client) List(ctx context.Context, req *generated.GetUsersRequest) (*generated.GetUsersResponseV2, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &generated.GetUsersRequest{}
	}
	result, err := cl.GetUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("users: list: %w", err)
	}
	return result, nil
}

// Get returns a single user by ID.
func (c *Client) Get(ctx context.Context, id string) (*generated.UserResponseV2, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	result, err := cl.GetUser0(ctx, &generated.GetUser0Request{
		ID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("users: get %s: %w", id, err)
	}
	return result, nil
}

// Create creates a new user.
func (c *Client) Create(ctx context.Context, req *generated.CreateUserRequest) (*generated.UserResponseV2, error) {
	cl, err := c.newGenerated(ctx)
	if err != nil {
		return nil, err
	}
	result, err := cl.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("users: create: %w", err)
	}
	return result, nil
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
