package client

import "context"

// TokenExchangeRequest contains the values sent to IBM Verify.
// Deprecated: Use client.New() + Client.Token.Exchange() instead.
type TokenExchangeRequest struct {
	TenantURL        string
	ClientID         string
	ClientSecret     string
	SubjectToken     string
	SubjectTokenType string
}

// TokenExchangeResponse represents the response from IBM Verify.
// Deprecated: Use ExchangeResult instead.
type TokenExchangeResponse = ExchangeResult

// ExchangeToken sends a JWT to IBM Verify and returns an access token.
//
// Deprecated: Use client.New() + Client.Token.Exchange() instead.
func ExchangeToken(
	ctx context.Context,
	req TokenExchangeRequest,
) (*ExchangeResult, error) {
	c, err := New(req.TenantURL,
		WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.Exchange(ctx, req.SubjectToken, req.SubjectTokenType)
}
