package client

import "context"

// SignerCertRequest contains the values needed to import a signer certificate.
// Deprecated: Use client.New() + Client.Certs.Import() instead.
type SignerCertRequest struct {
	TenantURL      string
	AccessToken    string
	CertificatePEM string
	Label          string
}

// SignerCertResponse represents a signer certificate returned by IBM Verify.
// Deprecated: Use CertResult instead.
type SignerCertResponse = CertResult

// ImportSignerCert uploads a signer certificate to IBM Verify.
//
// Deprecated: Use client.New() + Client.Certs.Import() instead.
// Note: this shim creates a temporary Client using the AccessToken directly —
// the new Client.Certs.Import() obtains its own token via ClientCredentials.
func ImportSignerCert(
	ctx context.Context,
	req SignerCertRequest,
) error {
	c, err := New(req.TenantURL)
	if err != nil {
		return err
	}
	// Shim: inject the pre-obtained token directly into CertsClient.
	return c.Certs.importWithToken(ctx, req.Label, req.CertificatePEM, req.AccessToken)
}

// GetSignerCert retrieves a signer certificate by label.
//
// Deprecated: Use client.New() + Client.Certs.Get() instead.
func GetSignerCert(
	ctx context.Context,
	tenantURL string,
	label string,
	accessToken string,
) (*CertResult, error) {
	c, err := New(tenantURL)
	if err != nil {
		return nil, err
	}
	return c.Certs.getWithToken(ctx, label, accessToken)
}

// DeleteSignerCert deletes a signer certificate by label.
//
// Deprecated: Use client.New() + Client.Certs.Delete() instead.
func DeleteSignerCert(
	ctx context.Context,
	tenantURL string,
	label string,
	accessToken string,
) error {
	c, err := New(tenantURL)
	if err != nil {
		return err
	}
	return c.Certs.deleteWithToken(ctx, label, accessToken)
}
