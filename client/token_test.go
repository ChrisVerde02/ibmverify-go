package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeVerify creates a test HTTP server handling all IBM Verify token endpoints.
func fakeVerify(t *testing.T, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for path, h := range handlers {
		mux.HandleFunc(path, h)
	}
	return httptest.NewServer(mux)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func TestTokenClientCredentials_success(t *testing.T) {
	srv := fakeVerify(t, map[string]http.HandlerFunc{
		"/v1.0/endpoint/default/token": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.FormValue("grant_type") != "client_credentials" {
				t.Errorf("wrong grant_type: %s", r.FormValue("grant_type"))
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"access_token": "test-token",
				"expires_in":   3600,
				"token_type":   "Bearer",
				"scope":        "openid",
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("client-id", "client-secret"))
	result, err := c.Token.ClientCredentials(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccessToken != "test-token" {
		t.Errorf("expected 'test-token', got %q", result.AccessToken)
	}
	if result.ExpiresIn != 3600 {
		t.Errorf("expected 3600, got %d", result.ExpiresIn)
	}
}

func TestTokenClientCredentials_missingClientID(t *testing.T) {
	c, _ := New("https://example.verify.ibm.com")
	_, err := c.Token.ClientCredentials(context.Background())
	if err == nil {
		t.Fatal("expected error for missing client ID")
	}
}

func TestTokenClientCredentials_http401(t *testing.T) {
	srv := fakeVerify(t, map[string]http.HandlerFunc{
		"/v1.0/endpoint/default/token": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"error":             "invalid_client",
				"error_description": "bad credentials",
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("bad-id", "bad-secret"))
	_, err := c.Token.ClientCredentials(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestTokenExchange_success(t *testing.T) {
	srv := fakeVerify(t, map[string]http.HandlerFunc{
		"/oauth2/token": func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			if r.FormValue("grant_type") != "urn:ietf:params:oauth:grant-type:token-exchange" {
				t.Errorf("wrong grant_type: %s", r.FormValue("grant_type"))
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"access_token":      "exchanged-token",
				"expires_in":        3599,
				"token_type":        "bearer",
				"issued_token_type": "urn:ietf:params:oauth:token-type:access_token",
				"grant_id":          "grant-abc",
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("sts-id", "sts-secret"))
	result, err := c.Token.Exchange(context.Background(),
		"subject-jwt",
		"urn:ietf:params:oauth:token-type:jwt",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AccessToken != "exchanged-token" {
		t.Errorf("expected 'exchanged-token', got %q", result.AccessToken)
	}
	if result.GrantID != "grant-abc" {
		t.Errorf("expected 'grant-abc', got %q", result.GrantID)
	}
}

func TestTokenExchange_emptySubjectToken(t *testing.T) {
	c, _ := New("https://example.verify.ibm.com", WithClientCredentials("id", "secret"))
	_, err := c.Token.Exchange(context.Background(), "", "urn:ietf:params:oauth:token-type:jwt")
	if err == nil {
		t.Fatal("expected error for empty subject token")
	}
}

func TestTokenIntrospect_success(t *testing.T) {
	srv := fakeVerify(t, map[string]http.HandlerFunc{
		"/oauth2/introspect": func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			if r.FormValue("token") == "" {
				t.Error("expected token in form, got empty")
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"active":   true,
				"sub":      "user123",
				"username": "testuser",
				"exp":      9999999999,
			})
		},
	})
	defer srv.Close()

	c, _ := New(srv.URL, WithClientCredentials("id", "secret"))
	result, err := c.Token.Introspect(context.Background(), "some-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Active {
		t.Error("expected active=true")
	}
	if result.Subject != "user123" {
		t.Errorf("expected 'user123', got %q", result.Subject)
	}
}
