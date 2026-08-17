package legacy

import (
	"context"

	"github.com/ChrisVerde02/ibmverify-go/client"
)

// TokenExchangeRequest contains the values sent to IBM Verify.
// Deprecated: Use client.New() + Client.Token.Exchange() instead.
type TokenExchangeRequest = client.TokenExchangeRequest

// TokenExchangeResponse represents the response from IBM Verify.
// Deprecated: Use client.ExchangeResult instead.
type TokenExchangeResponse = client.ExchangeResult

// ExchangeToken sends a JWT to IBM Verify and returns an access token.
//
// Deprecated: Use client.New() + Client.Token.Exchange() instead.
func ExchangeToken(
	ctx context.Context,
	req client.TokenExchangeRequest,
) (*client.ExchangeResult, error) {
	c, err := client.New(req.TenantURL,
		client.WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.Exchange(ctx, req.SubjectToken, req.SubjectTokenType)
}
