// Package legacy provides backwards-compatible shim functions for callers
// that used the original standalone function API before v1.5.0.
//
// All functions in this package are deprecated. Migrate to client.New() and
// the domain clients (Client.Token, Client.Certs) instead.
package legacy

import (
	"context"

	"github.ibm.com/Christian-Verderame/ibmverify-go/client"
)

// GetClientCredentialsToken fetches an access token using the OAuth 2.0
// client credentials grant.
//
// Deprecated: Use client.New() + Client.Token.ClientCredentials() instead.
func GetClientCredentialsToken(
	ctx context.Context,
	req client.ClientCredentialsRequest,
) (*client.ClientCredentialsResult, error) {
	c, err := client.New(req.TenantURL,
		client.WithClientCredentials(req.ClientID, req.ClientSecret),
	)
	if err != nil {
		return nil, err
	}
	return c.Token.ClientCredentials(ctx)
}
