package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SignerCertRequest contains the values needed to import a signer certificate.
type SignerCertRequest struct {
	TenantURL   string
	AccessToken string
	// CertificatePEM is the raw PEM string sent directly to IBM Verify.
	CertificatePEM string
	// Label is the friendly name / alias for the certificate in IBM Verify.
	// This must match the kid header used in the JWT (key_id in the provider).
	Label string
}

// SignerCertResponse represents a signer certificate returned by IBM Verify.
type SignerCertResponse struct {
	Label   string `json:"label"`
	Cert    string `json:"cert"`
	Subject string `json:"subjectDN"`
	Issuer  string `json:"issuerDN"`
}

// importSignerCertBody is the JSON body sent to POST /v1.0/signercert.
type importSignerCertBody struct {
	// Cert is the raw PEM string of the certificate.
	Cert  string `json:"cert"`
	Label string `json:"label"`
}

// ImportSignerCert uploads a signer certificate to IBM Verify.
// It calls POST /v1.0/signercert with the base64-encoded PEM and label.
// The access token must belong to a client with the manageCerts entitlement.
func ImportSignerCert(
	ctx context.Context,
	request SignerCertRequest,
) error {
	if strings.TrimSpace(request.TenantURL) == "" {
		return errors.New("tenant URL cannot be empty")
	}

	if strings.TrimSpace(request.AccessToken) == "" {
		return errors.New("access token cannot be empty")
	}

	if strings.TrimSpace(request.CertificatePEM) == "" {
		return errors.New("certificate PEM cannot be empty")
	}

	if strings.TrimSpace(request.Label) == "" {
		return errors.New("certificate label cannot be empty")
	}

	endpoint := strings.TrimRight(request.TenantURL, "/") + "/v1.0/signercert"

	body := importSignerCertBody{
		// IBM Verify expects the raw PEM string in the cert field.
		Cert: request.CertificatePEM,
		Label: request.Label,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal signer cert request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return fmt.Errorf("create signer cert request: %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+request.AccessToken)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("send signer cert request: %w", err)
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return fmt.Errorf("read signer cert response: %w", err)
	}

	// 201 Created is the success response per the API docs.
	if httpResponse.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"IBM Verify import signer cert failed with HTTP %d: %s",
			httpResponse.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

// GetSignerCert retrieves a signer certificate by label.
// It calls GET /v1.0/signercert/{label}.
// The access token must belong to a client with the readCerts entitlement.
func GetSignerCert(
	ctx context.Context,
	tenantURL string,
	label string,
	accessToken string,
) (*SignerCertResponse, error) {
	if strings.TrimSpace(tenantURL) == "" {
		return nil, errors.New("tenant URL cannot be empty")
	}

	if strings.TrimSpace(label) == "" {
		return nil, errors.New("certificate label cannot be empty")
	}

	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("access token cannot be empty")
	}

	endpoint := strings.TrimRight(tenantURL, "/") + "/v1.0/signercert/" + label

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create get signer cert request: %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
	httpRequest.Header.Set("Accept", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send get signer cert request: %w", err)
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("read get signer cert response: %w", err)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, nil // cert does not exist — not an error
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"IBM Verify get signer cert failed with HTTP %d: %s",
			httpResponse.StatusCode,
			string(responseBody),
		)
	}

	var result SignerCertResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("decode signer cert response: %w", err)
	}

	return &result, nil
}

// DeleteSignerCert deletes a signer certificate by label.
// It calls DELETE /v1.0/signercert/{label}.
// The access token must belong to a client with the manageCerts entitlement.
func DeleteSignerCert(
	ctx context.Context,
	tenantURL string,
	label string,
	accessToken string,
) error {
	if strings.TrimSpace(tenantURL) == "" {
		return errors.New("tenant URL cannot be empty")
	}

	if strings.TrimSpace(label) == "" {
		return errors.New("certificate label cannot be empty")
	}

	if strings.TrimSpace(accessToken) == "" {
		return errors.New("access token cannot be empty")
	}

	endpoint := strings.TrimRight(tenantURL, "/") + "/v1.0/signercert/" + label

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		endpoint,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create delete signer cert request: %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
	httpRequest.Header.Set("Accept", "application/json")

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("send delete signer cert request: %w", err)
	}
	defer httpResponse.Body.Close()

	// 204 No Content is the success response per the API docs.
	if httpResponse.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf(
			"IBM Verify delete signer cert failed with HTTP %d: %s",
			httpResponse.StatusCode,
			string(responseBody),
		)
	}

	return nil
}
