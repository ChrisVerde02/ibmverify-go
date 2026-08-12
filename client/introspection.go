package client

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
