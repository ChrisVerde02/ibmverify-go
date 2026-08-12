package client

import "context"

// ClientCredentialsRequest contains the values sent to IBM Verify.
// Deprecated: Use client.New() + Client.Token.ClientCredentials() instead.
type ClientCredentialsRequest struct {
	TenantURL    string
	ClientID     string
	ClientSecret string
}

// ClientCredentialsResponse represents the token returned by IBM Verify.
// Deprecated: Use ClientCredentialsResult instead.
type ClientCredentialsResponse = ClientCredentialsResult

// GetClientCredentialsToken fetches an access token using the OAuth 2.0
// client credentials grant.
//
// Deprecated: Use client.New() + Client.Token.ClientCredentials() instead.
func GetClientCredentialsToken(
	ctx context.Context,
	req ClientCredentialsRequest,
) (*ClientCredentialsResult, error) {
	c, err := New(req.TenantURL,
		WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.ClientCredentials(ctx)
}
