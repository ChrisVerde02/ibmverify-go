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
