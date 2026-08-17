package client

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError is returned by all SDK methods when IBM Verify responds with a
// non-2xx status code. It carries the HTTP status, IBM error code, and a
// human-readable message so callers can branch on typed errors instead of
// parsing strings.
//
// Usage:
//
//	var apiErr *client.APIError
//	if errors.As(err, &apiErr) {
//	    if apiErr.IsNotFound() { ... }
//	    if apiErr.IsAuth()     { ... }
//	    fmt.Println(apiErr.StatusCode, apiErr.Message)
//	}
type APIError struct {
	// StatusCode is the HTTP status returned by IBM Verify.
	StatusCode int
	// Code is the IBM-specific error code, e.g. "CSIAO5401E" or "invalid_client".
	// Empty if IBM Verify did not return a structured error body.
	Code string
	// Message is the human-readable error description.
	Message string
	// Endpoint is the URL path that produced the error.
	Endpoint string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("IBM Verify failed with HTTP %d: %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("IBM Verify failed with HTTP %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether the error is an HTTP 404.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsAuth reports whether the error is an authentication or authorisation failure.
func (e *APIError) IsAuth() bool {
	return e.StatusCode == http.StatusUnauthorized ||
		e.StatusCode == http.StatusForbidden ||
		e.Code == "invalid_client" ||
		e.Code == "unauthorized_client"
}

// IsRateLimit reports whether IBM Verify asked us to back off.
func (e *APIError) IsRateLimit() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServer reports whether the error is a 5xx server-side failure.
func (e *APIError) IsServer() bool {
	return e.StatusCode >= 500
}

// IsRetryable reports whether the request can safely be retried.
// Only 429 and 5xx transient errors are retryable.
func (e *APIError) IsRetryable() bool {
	return e.IsRateLimit() || e.IsServer()
}

// ErrNotFound is a sentinel returned by Get operations when IBM Verify
// responds with 404. Use errors.Is to test for it without caring about
// the full APIError detail.
var ErrNotFound = errors.New("not found")
