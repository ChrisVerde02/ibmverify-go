package legacy

import (
	"context"

	"github.com/ChrisVerde02/ibmverify-go/client"
)

// IntrospectionRequest contains the values sent to IBM Verify.
// Deprecated: Use client.New() + Client.Token.Introspect() instead.
type IntrospectionRequest = client.IntrospectionRequest

// IntrospectionResponse represents token metadata returned by IBM Verify.
// Deprecated: Use client.IntrospectResult instead.
type IntrospectionResponse = client.IntrospectResult

// IntrospectToken asks IBM Verify for metadata about an access token.
//
// Deprecated: Use client.New() + Client.Token.Introspect() instead.
func IntrospectToken(
	ctx context.Context,
	req client.IntrospectionRequest,
) (*client.IntrospectResult, error) {
	c, err := client.New(req.TenantURL,
		client.WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.Introspect(ctx, req.Token)
}
