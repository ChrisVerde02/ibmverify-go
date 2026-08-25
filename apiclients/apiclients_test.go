package apiclients_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChrisVerde02/ibmverify-go/apiclients"
	generated "github.com/ChrisVerde02/ibmverify-go/generated"
)

// fakeServer stubs the token endpoint and the /v1.0/apiclients endpoints.
func fakeServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// token endpoint
	mux.HandleFunc("/v1.0/endpoint/default/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "test-token",
			"expires_in":   3600,
		})
	})

	// list + create
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]any{
				map[string]any{"clientId": "cid-001", "clientName": "App One", "enabled": true},
				map[string]any{"clientId": "cid-002", "clientName": "App Two", "enabled": false},
			})
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"clientId":     "cid-new",
				"clientName":   "New Client",
				"clientSecret": "s3cr3t",
				"enabled":      true,
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// get + delete single client
	mux.HandleFunc("/v1.0/apiclients/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"clientId":   "cid-001",
				"clientName": "App One",
				"enabled":    true,
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewServer(mux)
}

func newClient(t *testing.T, tenantURL string) *apiclients.Client {
	t.Helper()
	return apiclients.New(tenantURL, func(ctx context.Context) (string, error) {
		return "test-token", nil
	})
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestList_returnsAllClients(t *testing.T) {
	srv := fakeServer(t)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	list, err := cl.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(list))
	}
	if list[0]["clientId"] != "cid-001" {
		t.Errorf("first client ID: got %v", list[0]["clientId"])
	}
	if list[1]["clientName"] != "App Two" {
		t.Errorf("second client name: got %v", list[1]["clientName"])
	}
}

func TestList_serverError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	_, err := cl.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestList_emptyArray(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	list, err := cl.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List empty: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet_returnsClient(t *testing.T) {
	srv := fakeServer(t)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	m, err := cl.Get(context.Background(), "cid-001")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if m["clientId"] != "cid-001" {
		t.Errorf("clientId: got %v", m["clientId"])
	}
	if m["clientName"] != "App One" {
		t.Errorf("clientName: got %v", m["clientName"])
	}
}

func TestGet_notFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"messageID":"CSIAK0001E","messageDescription":"Not found"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	_, err := cl.Get(context.Background(), "no-such-client")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate_returnsNewClient(t *testing.T) {
	srv := fakeServer(t)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	result, err := cl.Create(context.Background(), &generated.APIClientConfigRequest{
		ClientName:   "New Client",
		Entitlements: []string{"manageOAuthClients"},
		Enabled:      true,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if result["clientId"] != "cid-new" {
		t.Errorf("clientId: got %v", result["clientId"])
	}
	if result["clientSecret"] != "s3cr3t" {
		t.Errorf("clientSecret: got %v", result["clientSecret"])
	}
}

func TestCreate_sendsJSONBody(t *testing.T) {
	var received map[string]any
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"clientId": "x", "clientName": received["clientName"]})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	_, err := cl.Create(context.Background(), &generated.APIClientConfigRequest{
		ClientName:   "TestClient",
		Entitlements: []string{"manageOAuthClients"},
		Enabled:      true,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if received["clientName"] != "TestClient" {
		t.Errorf("request body clientName: got %v", received["clientName"])
	}
	if received["enabled"] != true {
		t.Errorf("request body enabled: got %v", received["enabled"])
	}
}

func TestCreate_fullResponseFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"clientId":     "cid-full",
			"clientName":   "Full Client",
			"clientSecret": "full-s3cr3t",
			"enabled":      true,
			"entitlements": []string{"manageOAuthClients", "manageCerts"},
			"description":  "test client",
			"ipFilterOp":   "WHITELIST",
			"ipFilters":    []string{"10.0.0.0/8"},
			"jwkUri":       "https://example.com/.well-known/jwks.json",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	result, err := cl.Create(context.Background(), &generated.APIClientConfigRequest{
		ClientName:   "Full Client",
		Entitlements: []string{"manageOAuthClients", "manageCerts"},
		Enabled:      true,
	})
	if err != nil {
		t.Fatalf("Create full: %v", err)
	}

	checks := map[string]any{
		"clientId":    "cid-full",
		"clientName":  "Full Client",
		"clientSecret": "full-s3cr3t",
		"description": "test client",
		"ipFilterOp":  "WHITELIST",
		"jwkUri":      "https://example.com/.well-known/jwks.json",
	}
	for k, want := range checks {
		if result[k] != want {
			t.Errorf("field %q: want %v, got %v", k, want, result[k])
		}
	}

	// entitlements comes back as []interface{} after JSON round-trip
	ents, ok := result["entitlements"].([]interface{})
	if !ok {
		t.Fatalf("entitlements: expected []interface{}, got %T", result["entitlements"])
	}
	if len(ents) != 2 {
		t.Errorf("entitlements: want 2, got %d", len(ents))
	}
	if ents[0] != "manageOAuthClients" {
		t.Errorf("entitlements[0]: want manageOAuthClients, got %v", ents[0])
	}

	// ipFilters same
	filters, ok := result["ipFilters"].([]interface{})
	if !ok {
		t.Fatalf("ipFilters: expected []interface{}, got %T", result["ipFilters"])
	}
	if len(filters) != 1 || filters[0] != "10.0.0.0/8" {
		t.Errorf("ipFilters: want [10.0.0.0/8], got %v", filters)
	}
}

func TestGet_fullResponseFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"clientId":     "cid-full",
			"clientName":   "Full Client",
			"enabled":      false,
			"entitlements": []string{"manageOAuthClients"},
			"description":  "read back",
			"ipFilterOp":   "BLACKLIST",
			"ipFilters":    []string{"192.168.0.0/16", "172.16.0.0/12"},
			"jwkUri":       "https://keys.example.com/jwks",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	m, err := cl.Get(context.Background(), "cid-full")
	if err != nil {
		t.Fatalf("Get full: %v", err)
	}

	if m["enabled"] != false {
		t.Errorf("enabled: want false, got %v", m["enabled"])
	}
	if m["description"] != "read back" {
		t.Errorf("description: want 'read back', got %v", m["description"])
	}
	if m["ipFilterOp"] != "BLACKLIST" {
		t.Errorf("ipFilterOp: want BLACKLIST, got %v", m["ipFilterOp"])
	}
	if m["jwkUri"] != "https://keys.example.com/jwks" {
		t.Errorf("jwkUri: got %v", m["jwkUri"])
	}

	filters, ok := m["ipFilters"].([]interface{})
	if !ok {
		t.Fatalf("ipFilters: expected []interface{}, got %T", m["ipFilters"])
	}
	if len(filters) != 2 {
		t.Errorf("ipFilters: want 2 entries, got %d", len(filters))
	}

	// clientSecret must NOT be present on Get (IBM never returns it after creation)
	if _, present := m["clientSecret"]; present {
		t.Error("clientSecret must not be present in Get response")
	}
}

func TestList_fullResponseFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]any{
			map[string]any{
				"clientId":     "cid-a",
				"clientName":   "Client A",
				"enabled":      true,
				"entitlements": []string{"manageOAuthClients"},
				"description":  "first",
			},
			map[string]any{
				"clientId":     "cid-b",
				"clientName":   "Client B",
				"enabled":      false,
				"entitlements": []string{"manageCerts"},
				"description":  "second",
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	list, err := cl.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List full: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("want 2, got %d", len(list))
	}
	if list[0]["description"] != "first" {
		t.Errorf("list[0] description: got %v", list[0]["description"])
	}
	if list[1]["description"] != "second" {
		t.Errorf("list[1] description: got %v", list[1]["description"])
	}

	// clientSecret must NOT be present in list items
	for i, item := range list {
		if _, present := item["clientSecret"]; present {
			t.Errorf("list[%d]: clientSecret must not be present in List response", i)
		}
	}

	ents, ok := list[0]["entitlements"].([]interface{})
	if !ok {
		t.Fatalf("list[0] entitlements: expected []interface{}, got %T", list[0]["entitlements"])
	}
	if ents[0] != "manageOAuthClients" {
		t.Errorf("list[0] entitlements[0]: got %v", ents[0])
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_ok(t *testing.T) {
	// Delete uses the Fern generated client, which needs a token endpoint too.
	srv := fakeServer(t)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	if err := cl.Delete(context.Background(), "cid-001"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestDelete_notFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1.0/apiclients/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"messageID":"CSIAK0001E","messageDescription":"Not found"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cl := newClient(t, srv.URL)
	err := cl.Delete(context.Background(), "no-such-client")
	if err == nil {
		t.Fatal("expected error for 404 delete, got nil")
	}
}
