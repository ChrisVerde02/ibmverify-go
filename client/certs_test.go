package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fakeCertVerify(t *testing.T, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	// All Certs methods need a client credentials token first
	mux.HandleFunc("/v1.0/endpoint/default/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "mgmt-token",
			"expires_in":   3600,
		})
	})
	for path, h := range handlers {
		mux.HandleFunc(path, h)
	}
	return httptest.NewServer(mux)
}

func TestCertsImport_success(t *testing.T) {
	srv := fakeCertVerify(t, map[string]http.HandlerFunc{
		"/v1.0/signercert": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			auth := r.Header.Get("Authorization")
			if auth != "Bearer mgmt-token" {
				t.Errorf("expected bearer token, got %q", auth)
			}
			w.WriteHeader(http.StatusCreated)
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("cm-id", "cm-secret"))
	err := c.Certs.Import(context.Background(), "testlabel", "-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCertsImport_emptyLabel(t *testing.T) {
	c, _ := New("https://example.verify.ibm.com", WithClientCredentials("id", "secret"))
	err := c.Certs.Import(context.Background(), "", "pem")
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestCertsImport_emptyPEM(t *testing.T) {
	c, _ := New("https://example.verify.ibm.com", WithClientCredentials("id", "secret"))
	err := c.Certs.Import(context.Background(), "label", "")
	if err == nil {
		t.Fatal("expected error for empty PEM")
	}
}

func TestCertsGet_found(t *testing.T) {
	srv := fakeCertVerify(t, map[string]http.HandlerFunc{
		"/v1.0/signercert/testlabel": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"label":    "testlabel",
				"cert":     "-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----",
				"subjectDN": "CN=testlabel",
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("cm-id", "cm-secret"))
	result, err := c.Certs.Get(context.Background(), "testlabel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected cert result, got nil")
	}
	if result.Label != "testlabel" {
		t.Errorf("expected label 'testlabel', got %q", result.Label)
	}
}

func TestCertsGet_notFound(t *testing.T) {
	srv := fakeCertVerify(t, map[string]http.HandlerFunc{
		"/v1.0/signercert/missing": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("cm-id", "cm-secret"))
	result, err := c.Certs.Get(context.Background(), "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for 404, got %+v", result)
	}
}

func TestCertsDelete_success(t *testing.T) {
	srv := fakeCertVerify(t, map[string]http.HandlerFunc{
		"/v1.0/signercert/testlabel": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("cm-id", "cm-secret"))
	err := c.Certs.Delete(context.Background(), "testlabel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCertsDelete_emptyLabel(t *testing.T) {
	c, _ := New("https://example.verify.ibm.com", WithClientCredentials("id", "secret"))
	err := c.Certs.Delete(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestCertsGet_normalisesBase64DER(t *testing.T) {
	// IBM Verify returns cert as raw base64 DER, not PEM.
	// The SDK must normalise it to PEM so callers can compare with stored PEM.
	srv := fakeCertVerify(t, map[string]http.HandlerFunc{
		"/v1.0/signercert/testlabel": func(w http.ResponseWriter, r *http.Request) {
			// Return the cert as raw base64 DER (what IBM Verify actually sends)
			w.Header().Set("Content-Type", "application/json")
			// A minimal valid DER-encoded cert encoded as base64 would be complex,
			// so we test the normalisation path using a known base64 string that
			// does NOT start with "-----". The result must start with the PEM header.
			json.NewEncoder(w).Encode(map[string]any{
				"label": "testlabel",
				// This is the base64 of "fake" — not a real cert but tests the
				// normalisation code path (base64 → PEM wrapping).
				"cert":     "ZmFrZQ==",
				"subjectDN": "CN=testlabel",
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("cm-id", "cm-secret"))
	result, err := c.Certs.Get(context.Background(), "testlabel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected cert result, got nil")
	}
	if !strings.HasPrefix(result.Cert, "-----BEGIN CERTIFICATE-----") {
		t.Errorf("expected PEM-normalised cert, got: %q", result.Cert[:min(len(result.Cert), 60)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
