package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestGenerateSelfSignedCertificate_success(t *testing.T) {
	result, err := GenerateSelfSignedCertificate(CertificateRequest{
		CommonName:   "TestCert",
		Organization: "IBM",
		Country:      "US",
		ValidityDays: 365,
		KeySize:      2048,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Cert PEM must be parseable
	block, _ := pem.Decode([]byte(result.CertificatePEM))
	if block == nil {
		t.Fatal("certificate PEM did not decode")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("certificate parse failed: %v", err)
	}
	if cert.Subject.CommonName != "TestCert" {
		t.Errorf("expected CN=TestCert, got %q", cert.Subject.CommonName)
	}

	// Private key PEM must be parseable
	keyBlock, _ := pem.Decode([]byte(result.PrivateKeyPEM))
	if keyBlock == nil {
		t.Fatal("private key PEM did not decode")
	}
}

func TestGenerateSelfSignedCertificate_noExtKeyUsageClientAuth(t *testing.T) {
	result, err := GenerateSelfSignedCertificate(CertificateRequest{
		CommonName:   "SignerCert",
		Organization: "IBM",
		Country:      "US",
		ValidityDays: 1,
		KeySize:      2048,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	block, _ := pem.Decode([]byte(result.CertificatePEM))
	cert, _ := x509.ParseCertificate(block.Bytes)

	for _, eku := range cert.ExtKeyUsage {
		if eku == x509.ExtKeyUsageClientAuth {
			t.Error("cert must NOT have ExtKeyUsageClientAuth — IBM Verify rejects it")
		}
	}
}

func TestGenerateSelfSignedCertificate_invalidKeySize(t *testing.T) {
	_, err := GenerateSelfSignedCertificate(CertificateRequest{
		CommonName:   "Test",
		Organization: "IBM",
		Country:      "US",
		ValidityDays: 1,
		KeySize:      1024, // not allowed
	})
	if err == nil {
		t.Fatal("expected error for key size 1024")
	}
}

func TestGenerateSelfSignedCertificate_missingFields(t *testing.T) {
	tests := []struct {
		name string
		req  CertificateRequest
	}{
		{"no CN", CertificateRequest{Organization: "IBM", Country: "US", ValidityDays: 1, KeySize: 2048}},
		{"no org", CertificateRequest{CommonName: "T", Country: "US", ValidityDays: 1, KeySize: 2048}},
		{"no country", CertificateRequest{CommonName: "T", Organization: "IBM", ValidityDays: 1, KeySize: 2048}},
		{"zero days", CertificateRequest{CommonName: "T", Organization: "IBM", Country: "US", KeySize: 2048}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateSelfSignedCertificate(tt.req)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
