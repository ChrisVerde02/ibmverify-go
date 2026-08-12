package client

import (
	"net/http"
	"testing"
	"time"
)

func TestNew_emptyTenantURL(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty tenant URL, got nil")
	}
}

func TestNew_trailingSlashStripped(t *testing.T) {
	c, err := New("https://example.verify.ibm.com/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.TenantURL() != "https://example.verify.ibm.com" {
		t.Errorf("expected trailing slash stripped, got %q", c.TenantURL())
	}
}

func TestWithTimeout(t *testing.T) {
	c, err := New("https://example.verify.ibm.com",
		WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", c.httpClient.Timeout)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 99 * time.Second}
	c, err := New("https://example.verify.ibm.com", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient != custom {
		t.Error("expected custom HTTP client to be set")
	}
}

func TestDomainClientsAttached(t *testing.T) {
	c, err := New("https://example.verify.ibm.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Token == nil {
		t.Error("expected Token domain client to be non-nil")
	}
	if c.Certs == nil {
		t.Error("expected Certs domain client to be non-nil")
	}
}
