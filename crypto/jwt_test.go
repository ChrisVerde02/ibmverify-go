package crypto

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateJTI_returnsUUID(t *testing.T) {
	jti, err := GenerateJTI()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	parts := strings.Split(jti, "-")
	if len(parts) != 5 {
		t.Errorf("expected UUID with 5 parts, got %q", jti)
	}
}

func TestGenerateJTI_unique(t *testing.T) {
	a, _ := GenerateJTI()
	b, _ := GenerateJTI()
	if a == b {
		t.Error("expected unique JTIs, got identical values")
	}
}

func makeKey(t *testing.T) string {
	t.Helper()
	result, err := GenerateSelfSignedCertificate(CertificateRequest{
		CommonName:   "test",
		Organization: "IBM",
		Country:      "US",
		ValidityDays: 1,
		KeySize:      2048,
	})
	if err != nil {
		t.Fatalf("generate cert: %v", err)
	}
	return result.PrivateKeyPEM
}

func TestGenerateSignedJWT_success(t *testing.T) {
	key := makeKey(t)
	jti, _ := GenerateJTI()

	result, err := GenerateSignedJWT(JWTRequest{
		Issuer:        "https://test.ibm.com",
		Subject:       "testuser",
		KeyID:         "testkey",
		JWTID:         jti,
		PrivateKeyPEM: key,
		ExpiresIn:     5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token == "" {
		t.Error("expected non-empty token")
	}
	// JWT has 3 base64 parts separated by "."
	if len(strings.Split(result.Token, ".")) != 3 {
		t.Errorf("expected 3-part JWT, got %q", result.Token)
	}
	if result.IssuedAt == 0 {
		t.Error("expected non-zero IssuedAt")
	}
	if result.ExpiresAt <= result.IssuedAt {
		t.Error("expected ExpiresAt > IssuedAt")
	}
}

func TestGenerateSignedJWT_missingFields(t *testing.T) {
	key := makeKey(t)
	jti, _ := GenerateJTI()

	tests := []struct {
		name string
		req  JWTRequest
	}{
		{"empty issuer", JWTRequest{Subject: "u", KeyID: "k", JWTID: jti, PrivateKeyPEM: key, ExpiresIn: time.Minute}},
		{"empty subject", JWTRequest{Issuer: "https://i", KeyID: "k", JWTID: jti, PrivateKeyPEM: key, ExpiresIn: time.Minute}},
		{"empty key id", JWTRequest{Issuer: "https://i", Subject: "u", JWTID: jti, PrivateKeyPEM: key, ExpiresIn: time.Minute}},
		{"empty jti", JWTRequest{Issuer: "https://i", Subject: "u", KeyID: "k", PrivateKeyPEM: key, ExpiresIn: time.Minute}},
		{"zero expires", JWTRequest{Issuer: "https://i", Subject: "u", KeyID: "k", JWTID: jti, PrivateKeyPEM: key}},
		{"bad key", JWTRequest{Issuer: "https://i", Subject: "u", KeyID: "k", JWTID: jti, PrivateKeyPEM: "notakey", ExpiresIn: time.Minute}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateSignedJWT(tt.req)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
