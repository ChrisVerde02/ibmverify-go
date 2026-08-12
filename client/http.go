package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to extract an OAuth error description first.
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if jsonErr := json.Unmarshal(body, &oauthErr); jsonErr == nil && oauthErr.Error != "" {
			return nil, fmt.Errorf("IBM Verify failed with HTTP %d: %s: %s",
				resp.StatusCode, oauthErr.Error, oauthErr.ErrorDescription)
		}
		return nil, fmt.Errorf("IBM Verify failed with HTTP %d: %s",
			resp.StatusCode, string(body))
	}

	return body, nil
}
