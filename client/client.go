// Package client provides a Client for the IBM Verify API.
//
// Usage:
//
//	c, err := client.New("https://example.verify.ibm.com",
//	    client.WithClientCredentials("client-id", "client-secret"),
//	)
//	token, err := c.Token.ClientCredentials(ctx)
//	cert,  err := c.Certs.Get(ctx, "demotokensigner")
package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"
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
