package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

// Import uploads a PEM certificate to IBM Verify as a signer certificate.
// label must match the kid header used in JWTs during token exchange.
//
//	POST /v1.0/signercert
func (cr *CertsClient) Import(ctx context.Context, label, certificatePEM string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("import cert: label cannot be empty")
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("import cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// Get retrieves a signer certificate by label.
// Returns nil, nil when the certificate does not exist (HTTP 404).
//
//	GET /v1.0/signercert/{label}
func (cr *CertsClient) Get(ctx context.Context, label string) (*CertResult, error) {
	if strings.TrimSpace(label) == "" {
		return nil, fmt.Errorf("get cert: label cannot be empty")
	}

	token, err := cr.c.Token.ClientCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cert: get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cr.c.endpoint("/v1.0/signercert/"+label), nil)
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

	body, _ := io.ReadAll(resp.Body)

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
	return &result, nil
}

// Delete removes a signer certificate from IBM Verify by label.
//
//	DELETE /v1.0/signercert/{label}
func (cr *CertsClient) Delete(ctx context.Context, label string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("delete cert: label cannot be empty")
	}

	token, err := cr.c.Token.ClientCredentials(ctx)
	if err != nil {
		return fmt.Errorf("delete cert: get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		cr.c.endpoint("/v1.0/signercert/"+label), nil)
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// importWithToken is used by the backwards-compat shim when a pre-obtained
// token is passed directly (legacy ImportSignerCert behaviour).
func (cr *CertsClient) importWithToken(ctx context.Context, label, certificatePEM, accessToken string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("import cert: label cannot be empty")
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("import cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// getWithToken is used by the backwards-compat shim.
func (cr *CertsClient) getWithToken(ctx context.Context, label, accessToken string) (*CertResult, error) {
	if strings.TrimSpace(label) == "" {
		return nil, fmt.Errorf("get cert: label cannot be empty")
	}
	if strings.TrimSpace(accessToken) == "" {
		return nil, fmt.Errorf("get cert: access token cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cr.c.endpoint("/v1.0/signercert/"+label), nil)
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

	body, _ := io.ReadAll(resp.Body)

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
	return &result, nil
}

// deleteWithToken is used by the backwards-compat shim.
func (cr *CertsClient) deleteWithToken(ctx context.Context, label, accessToken string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("delete cert: label cannot be empty")
	}
	if strings.TrimSpace(accessToken) == "" {
		return fmt.Errorf("delete cert: access token cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		cr.c.endpoint("/v1.0/signercert/"+label), nil)
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete cert: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
