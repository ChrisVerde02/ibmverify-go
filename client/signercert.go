package client

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
