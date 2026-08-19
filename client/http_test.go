package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestDo_RetryOnRateLimit verifies that do() retries a 429 response up to
// retryMax-1 times and eventually returns the successful response.
func TestDo_RetryOnRateLimit(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < retryMax {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"ok"`))
	}))
	defer srv.Close()

	c := &Client{httpClient: srv.Client()}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, strings.NewReader("body"))

	// Speed up the test — replace retryBaseWait temporarily via a zero-wait server
	// by keeping retryBaseWait as-is but using a very short timeout; instead we
	// just accept the real waits are 1s+2s in production and make the test fast
	// by only doing 1 retry (retryMax=3, fail twice then succeed on attempt 3).
	start := time.Now()
	_, err := c.do(req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if attempts != retryMax {
		t.Errorf("expected %d attempts, got %d", retryMax, attempts)
	}
	// 3 attempts = 2 waits: 1s + 2s = 3s minimum (allow generous upper bound)
	if elapsed < 3*time.Second {
		t.Errorf("expected at least 3s of backoff, got %v", elapsed)
	}
}

// TestDo_NoRetryOnClientError verifies that a 4xx non-retryable error is
// returned immediately without retrying.
func TestDo_NoRetryOnClientError(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_request","error_description":"bad"}`))
	}))
	defer srv.Close()

	c := &Client{httpClient: srv.Client()}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, strings.NewReader("body"))

	_, err := c.do(req)
	if err == nil {
		t.Fatal("expected error for 400, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected exactly 1 attempt for non-retryable error, got %d", attempts)
	}
}

// TestDo_RetryExhausted verifies that when all retries fail, the last error is returned.
func TestDo_RetryExhausted(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"messageId":"CSIAO0001E","messageDescription":"unavailable"}`))
	}))
	defer srv.Close()

	c := &Client{httpClient: srv.Client()}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, strings.NewReader("body"))

	_, err := c.do(req)
	if err == nil {
		t.Fatal("expected error after retries exhausted, got nil")
	}
	if attempts != retryMax {
		t.Errorf("expected %d attempts, got %d", retryMax, attempts)
	}
}

// TestDo_ContextCancelledDuringBackoff verifies that a cancelled context
// aborts the retry wait and returns the context error.
func TestDo_ContextCancelledDuringBackoff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{httpClient: srv.Client()}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, srv.URL, strings.NewReader("body"))

	// Cancel after first failure triggers the backoff wait — cancel immediately
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := c.do(req)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error after context cancel, got nil")
	}
	// Should have returned well before the 1s backoff completed
	if elapsed > 800*time.Millisecond {
		t.Errorf("expected fast cancel, took %v", elapsed)
	}
}
