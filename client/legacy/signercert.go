package legacy

import (
	"context"

	"github.ibm.com/Christian-Verderame/ibmverify-go/client"
)

// SignerCertRequest contains the values needed to import a signer certificate.
// Deprecated: Use client.New() + Client.Certs.Import() instead.
type SignerCertRequest = client.SignerCertRequest

// SignerCertResponse represents a signer certificate returned by IBM Verify.
// Deprecated: Use client.CertResult instead.
type SignerCertResponse = client.CertResult

// ImportSignerCert uploads a signer certificate to IBM Verify.
//
// Deprecated: Use client.New() + Client.Certs.Import() instead.
func ImportSignerCert(
	ctx context.Context,
	req client.SignerCertRequest,
) error {
	c, err := client.New(req.TenantURL)
	if err != nil {
		return err
	}
	return c.Certs.ImportWithToken(ctx, req.Label, req.CertificatePEM, req.AccessToken)
}

// GetSignerCert retrieves a signer certificate by label.
//
// Deprecated: Use client.New() + Client.Certs.Get() instead.
func GetSignerCert(
	ctx context.Context,
	tenantURL string,
	label string,
	accessToken string,
) (*client.CertResult, error) {
	c, err := client.New(tenantURL)
	if err != nil {
		return nil, err
	}
	return c.Certs.GetWithToken(ctx, label, accessToken)
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
	c, err := client.New(tenantURL)
	if err != nil {
		return err
	}
	return c.Certs.DeleteWithToken(ctx, label, accessToken)
}
