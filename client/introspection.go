package client

import "context"

// IntrospectionRequest contains the values sent to IBM Verify.
// Deprecated: Use client.New() + Client.Token.Introspect() instead.
type IntrospectionRequest struct {
	TenantURL    string
	ClientID     string
	ClientSecret string
	Token        string
}

// IntrospectionResponse represents token metadata returned by IBM Verify.
// Deprecated: Use IntrospectResult instead.
type IntrospectionResponse = IntrospectResult

// IntrospectToken asks IBM Verify for metadata about an access token.
//
// Deprecated: Use client.New() + Client.Token.Introspect() instead.
func IntrospectToken(
	ctx context.Context,
	req IntrospectionRequest,
) (*IntrospectResult, error) {
	c, err := New(req.TenantURL,
		WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.Introspect(ctx, req.Token)
}
