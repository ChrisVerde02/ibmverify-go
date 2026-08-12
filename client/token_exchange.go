package client

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
