package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	retryMax      = 3
	retryBaseWait = time.Second
)

// maxResponseBody caps how many bytes we read from any IBM Verify response.
// Protects against hostile or misbehaving servers OOM-ing the process.
const maxResponseBody = 10 * 1024 * 1024 // 10 MiB

// checkStatus converts a non-2xx HTTP response into a typed *APIError.
// It tries to parse IBM Verify's structured error body first; falls back
// to the raw body as the message. Callers pass the expected success code(s)
// so each endpoint documents its contract explicitly.
//
// For GET endpoints that treat 404 as "not found" (not an error), callers
// should check resp.StatusCode == http.StatusNotFound before calling
// checkStatus.
func checkStatus(resp *http.Response, body []byte, endpoint string, allowed ...int) error {
	for _, code := range allowed {
		if resp.StatusCode == code {
			return nil
		}
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Endpoint:   endpoint,
	}

	// Try IBM Verify structured JSON error body first.
	var structured struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		MessageID        string `json:"messageId"`
		MessageDesc      string `json:"messageDescription"`
	}
	if json.Unmarshal(body, &structured) == nil {
		switch {
		case structured.Error != "":
			apiErr.Code = structured.Error
			apiErr.Message = structured.ErrorDescription
		case structured.MessageID != "":
			apiErr.Code = structured.MessageID
			apiErr.Message = structured.MessageDesc
		default:
			apiErr.Message = string(body)
		}
	} else {
		apiErr.Message = string(body)
	}

	return apiErr
}

// postForm sends a POST with an application/x-www-form-urlencoded body.
// bearerToken is optional; pass "" to omit the Authorization header.
func (c *Client) postForm(ctx context.Context, path string, form url.Values, bearerToken string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.endpoint(path), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	return c.do(req)
}

// do executes an http.Request and returns the body, or an error for non-2xx.
// It returns a typed *APIError for all non-2xx responses so callers can use
// errors.As() to inspect status codes without string parsing.
//
// Retryable responses (429, 5xx) are retried up to retryMax times with
// exponential backoff. The request body is reset via req.GetBody between
// attempts; requests without a GetBody function are not retried.
func (c *Client) do(req *http.Request) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < retryMax; attempt++ {
		if attempt > 0 {
			// Reset body for retry — only possible when GetBody is set.
			if req.GetBody == nil {
				return nil, lastErr
			}
			newBody, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("reset request body: %w", err)
			}
			req.Body = newBody

			// Exponential backoff: 1s, 2s, … Abort if context is done.
			wait := retryBaseWait * (1 << (attempt - 1))
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(wait):
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("send request: %w", err)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBody))
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}

		if err := checkStatus(resp, body, req.URL.Path, http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusAccepted); err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && apiErr.IsRetryable() {
				lastErr = err
				continue
			}
			return nil, err
		}

		return body, nil
	}
	return nil, lastErr
}
