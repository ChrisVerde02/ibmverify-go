// Package client provides a Client for the IBM Verify API.
//
// Usage:
//
//	c, err := client.New("https://example.verify.ibm.com",
//	    client.WithClientCredentials("client-id", "client-secret"),
//	)
//	token, err := c.Token.ClientCredentials(ctx)
//	cert,  err := c.Certs.Get(ctx, "demotokensigner")
//	apps,  err := c.Apps.List(ctx, nil)
package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChrisVerde02/ibmverify-go/apiclients"
	"github.com/ChrisVerde02/ibmverify-go/apps"
	"github.com/ChrisVerde02/ibmverify-go/users"
)

const defaultTimeout = 30 * time.Second

// Client is the top-level IBM Verify API client.
// Create one with New() and reuse it across all operations.
// It is safe for concurrent use.
type Client struct {
	tenantURL    string
	clientID     string
	clientSecret string
	httpClient   *http.Client

	// Token provides methods for OAuth token operations.
	Token *TokenClient
	// Certs provides methods for signer certificate management.
	Certs *CertsClient
	// Apps provides methods for IBM Verify application management.
	Apps *apps.Client
	// Users provides methods for IBM Verify user management (SCIM v2).
	Users *users.Client
	// APIClients provides methods for dynamic client registration.
	APIClients *apiclients.Client
}

// Option configures a Client.
type Option func(*Client)

// WithClientCredentials sets the OAuth client ID and secret used by
// Token.ClientCredentials and as the default auth for domain clients.
func WithClientCredentials(clientID, clientSecret string) Option {
	return func(c *Client) {
		c.clientID = clientID
		c.clientSecret = clientSecret
	}
}

// WithHTTPClient replaces the default HTTP client.
// Use this to inject a custom transport, proxy, or test mock.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the HTTP request timeout. Default is 30 seconds.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient = &http.Client{Timeout: d}
	}
}

// New creates a new Client for the given IBM Verify tenant URL.
// tenantURL must be the base tenant URL, e.g. https://example.verify.ibm.com.
func New(tenantURL string, opts ...Option) (*Client, error) {
	tenantURL = strings.TrimRight(strings.TrimSpace(tenantURL), "/")
	if tenantURL == "" {
		return nil, fmt.Errorf("tenant URL cannot be empty")
	}

	c := &Client{
		tenantURL:  tenantURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Attach domain clients — they share the same http.Client and base URL.
	c.Token = &TokenClient{c: c}
	c.Certs = &CertsClient{c: c}

	// Generated domain clients use the token client for auth.
	getToken := func(ctx context.Context) (string, error) {
		t, err := c.Token.ClientCredentials(ctx)
		if err != nil {
			return "", err
		}
		return t.AccessToken, nil
	}
	c.Apps = apps.New(c.tenantURL, getToken)
	c.Users = users.New(c.tenantURL, getToken)
	c.APIClients = apiclients.New(c.tenantURL, getToken)

	return c, nil
}

// TenantURL returns the configured tenant base URL.
func (c *Client) TenantURL() string {
	return c.tenantURL
}

// endpoint builds a full URL by appending path to the tenant URL.
func (c *Client) endpoint(path string) string {
	return c.tenantURL + path
}
