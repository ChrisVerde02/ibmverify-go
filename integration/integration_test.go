//go:build integration

// Run with: go test -tags integration -v ./integration/ -timeout 30s
// Requires environment variables:
//   VERIFY_TENANT_URL    e.g. https://ChrisVerde02.verify.ibm.com
//   VERIFY_CLIENT_ID     your IBM Verify client ID
//   VERIFY_CLIENT_SECRET your IBM Verify client secret

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ChrisVerde02/ibmverify-go/client"
)

func newClient(t *testing.T) *client.Client {
	t.Helper()
	tenantURL := os.Getenv("VERIFY_TENANT_URL")
	clientID := os.Getenv("VERIFY_CLIENT_ID")
	clientSecret := os.Getenv("VERIFY_CLIENT_SECRET")
	// Fall back to STS client if no explicit client set
	if clientID == "" {
		clientID = os.Getenv("VERIFY_STS_CLIENT_ID")
		clientSecret = os.Getenv("VERIFY_STS_CLIENT_SECRET")
	}
	if tenantURL == "" || clientID == "" || clientSecret == "" {
		t.Skip("VERIFY_TENANT_URL / VERIFY_CLIENT_ID / VERIFY_CLIENT_SECRET not set")
	}
	c, err := client.New(tenantURL,
		client.WithClientCredentials(clientID, clientSecret),
	)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	return c
}

func TestToken_ClientCredentials(t *testing.T) {
	c := newClient(t)
	tok, err := c.Token.ClientCredentials(context.Background())
	if err != nil {
		t.Fatalf("ClientCredentials: %v", err)
	}
	if tok.AccessToken == "" {
		t.Fatal("got empty access token")
	}
	t.Logf("token type: %s, expires_in: %d", tok.TokenType, tok.ExpiresIn)
}

func TestApps_List(t *testing.T) {
	c := newClient(t)
	result, err := c.Apps.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Apps.List: %v", err)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("Apps response:\n%s", string(b))
	fmt.Printf("\n✓ Apps.List returned successfully\n")
}

func TestUsers_List(t *testing.T) {
	c := newClient(t)
	result, err := c.Users.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Users.List: %v", err)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("Users response:\n%s", string(b))
	fmt.Printf("\n✓ Users.List returned successfully\n")
}

func TestAPIClients_List(t *testing.T) {
	c := newClient(t)
	result, err := c.APIClients.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("APIClients.List: %v", err)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("APIClients response:\n%s", string(b))
	fmt.Printf("\n✓ APIClients.List returned successfully\n")
}

func TestCerts_Get(t *testing.T) {
	c := newClient(t)
	cert, err := c.Certs.Get(context.Background(), "demotokensigner")
	if err != nil {
		t.Logf("Certs.Get error (cert may not exist): %v", err)
		t.Skip("demotokensigner cert not found on tenant")
	}
	b, _ := json.MarshalIndent(cert, "", "  ")
	t.Logf("cert:\n%s", string(b))
	fmt.Printf("\n✓ Certs.Get returned successfully\n")
}
