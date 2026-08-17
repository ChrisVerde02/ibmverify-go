package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// labelRe validates signer certificate labels.
// Allows letters, digits, dots, hyphens, underscores — 1 to 128 characters.
// This also prevents path-traversal attacks ("..", "/", etc.) via the label
// in URLs like /v1.0/signercert/{label}.
var labelRe = regexp.MustCompile(`^[A-Za-z0-9._\-]{1,128}$`)

// validateLabel returns an error if label is empty or contains unsafe characters.
func validateLabel(op, label string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("%s: label cannot be empty", op)
	}
	if !labelRe.MatchString(label) {
		return fmt.Errorf("%s: label %q contains invalid characters — only letters, digits, dots, hyphens, and underscores are allowed (1–128 chars)", op, label)
	}
	return nil
}

// signerCertPath returns the escaped URL path for a labelled signer cert.
func signerCertPath(label string) string {
	return "/v1.0/signercert/" + url.PathEscape(label)
}

// CertsClient provides signer certificate management against IBM Verify.
// Access it via Client.Certs.
//
// All methods obtain a client credentials token automatically using the
// clientID and clientSecret configured on the parent Client.
type CertsClient struct {
	c *Client
}

// CertResult represents a signer certificate returned by IBM Verify.
type CertResult struct {
	Label   string `json:"label"`
	Cert    string `json:"cert"`
	Subject string `json:"subjectDN"`
	Issuer  string `json:"issuerDN"`
}


// normaliseCertPEM converts an IBM Verify cert field to a PEM string.
// IBM Verify returns the cert as raw base64-encoded DER. If the value is
// already a valid PEM block it is returned unchanged.
func normaliseCertPEM(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "-----") {
		// Already PEM — return as-is.
		return raw
	}
	// Strip any whitespace/newlines that may be embedded in the base64.
	raw = strings.ReplaceAll(raw, "\n", "")
	raw = strings.ReplaceAll(raw, " ", "")
	der, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		// Not valid base64 — return unchanged and let callers handle it.
		return raw
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

// Import uploads a PEM certificate to IBM Verify as a signer certificate.
// label must match the kid header used in JWTs during token exchange.
//
//	POST /v1.0/signercert
func (cr *CertsClient) Import(ctx context.Context, label, certificatePEM string) error {
	if err := validateLabel("import cert", label); err != nil {
		return err
	}
	if strings.TrimSpace(certificatePEM) == "" {
		return fmt.Errorf("import cert: certificate PEM cannot be empty")
	}

	token, err := cr.c.Token.ClientCredentials(ctx)
	if err != nil {
		return fmt.Errorf("import cert: get access token: %w", err)
	}

	bodyBytes, err := json.Marshal(struct {
		Cert  string `json:"cert"`
		Label string `json:"label"`
	}{Cert: certificatePEM, Label: label})
	if err != nil {
		return fmt.Errorf("import cert: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		cr.c.endpoint("/v1.0/signercert"), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("import cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("import cert: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
		return fmt.Errorf("import cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// Get retrieves a signer certificate by label.
// Returns nil, nil when the certificate does not exist (HTTP 404).
//
//	GET /v1.0/signercert/{label}
func (cr *CertsClient) Get(ctx context.Context, label string) (*CertResult, error) {
	if err := validateLabel("get cert", label); err != nil {
		return nil, err
	}

	token, err := cr.c.Token.ClientCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cert: get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cr.c.endpoint(signerCertPath(label)), nil)
	if err != nil {
		return nil, fmt.Errorf("get cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get cert: send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // cert does not exist — not an error
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cert: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result CertResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("get cert: decode response: %w", err)
	}
	result.Cert = normaliseCertPEM(result.Cert)
	return &result, nil
}

// Delete removes a signer certificate from IBM Verify by label.
//
//	DELETE /v1.0/signercert/{label}
func (cr *CertsClient) Delete(ctx context.Context, label string) error {
	if err := validateLabel("delete cert", label); err != nil {
		return err
	}

	token, err := cr.c.Token.ClientCredentials(ctx)
	if err != nil {
		return fmt.Errorf("delete cert: get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		cr.c.endpoint(signerCertPath(label)), nil)
	if err != nil {
		return fmt.Errorf("delete cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete cert: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
		return fmt.Errorf("delete cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// ImportWithToken is used by the backwards-compat shim when a pre-obtained
// token is passed directly (legacy ImportSignerCert behaviour).
func (cr *CertsClient) ImportWithToken(ctx context.Context, label, certificatePEM, accessToken string) error {
	if err := validateLabel("import cert", label); err != nil {
		return err
	}
	if strings.TrimSpace(certificatePEM) == "" {
		return fmt.Errorf("import cert: certificate PEM cannot be empty")
	}
	if strings.TrimSpace(accessToken) == "" {
		return fmt.Errorf("import cert: access token cannot be empty")
	}

	bodyBytes, err := json.Marshal(struct {
		Cert  string `json:"cert"`
		Label string `json:"label"`
	}{Cert: certificatePEM, Label: label})
	if err != nil {
		return fmt.Errorf("import cert: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		cr.c.endpoint("/v1.0/signercert"), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("import cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("import cert: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
		return fmt.Errorf("import cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// GetWithToken is used by the backwards-compat shim.
func (cr *CertsClient) GetWithToken(ctx context.Context, label, accessToken string) (*CertResult, error) {
	if err := validateLabel("get cert", label); err != nil {
		return nil, err
	}
	if strings.TrimSpace(accessToken) == "" {
		return nil, fmt.Errorf("get cert: access token cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cr.c.endpoint(signerCertPath(label)), nil)
	if err != nil {
		return nil, fmt.Errorf("get cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get cert: send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cert: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result CertResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("get cert: decode response: %w", err)
	}
	result.Cert = normaliseCertPEM(result.Cert)
	return &result, nil
}

// DeleteWithToken is used by the backwards-compat shim.
func (cr *CertsClient) DeleteWithToken(ctx context.Context, label, accessToken string) error {
	if err := validateLabel("delete cert", label); err != nil {
		return err
	}
	if strings.TrimSpace(accessToken) == "" {
		return fmt.Errorf("delete cert: access token cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		cr.c.endpoint(signerCertPath(label)), nil)
	if err != nil {
		return fmt.Errorf("delete cert: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := cr.c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete cert: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
		return fmt.Errorf("delete cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
