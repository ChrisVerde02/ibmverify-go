package client

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
