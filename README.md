# ibmverify-go

Go SDK for [IBM Verify](https://www.ibm.com/products/verify-identity). Provides typed clients for token management, certificate management, application management, user management, and dynamic client registration — used by [`ibmverify-cli`](https://github.com/ChrisVerde02/ibmverify-cli) and [`terraform-provider-verify`](https://github.com/ChrisVerde02/terraform-provider-verify).

## Installation

```bash
go get github.com/ChrisVerde02/ibmverify-go
```

## Architecture

```
specs/openapi_corrected.yaml   ← canonical IBM Verify spec (385 paths)
        │
        │  fern generate
        ▼
generated/                     ← Fern output — DO NOT EDIT
  core/                        retry engine, APIError, HTTP types
  internal/                    Caller, Retrier (exponential backoff, jitter)
  applicationaccess/           generated Applications client
  usersmanagementversion20/    generated Users (SCIM v2) client
  apiclients/                  generated API Clients client
  types.go                     all IBM Verify request/response structs
        │
        │  thin handwritten wrappers (absorb IBM spec/reality mismatches)
        ▼
apps/           List/Get use raw HTTP (IBM returns applicationState as bool, spec says string)
users/          Accept: application/scim+json required; filter support
apiclients/     List/Get/Create use raw HTTP (Fern raw client returns Body: nil for DCR)
client/         top-level Client — Token, Certs, Apps, Users, APIClients
crypto/         local JWT signing and certificate generation
        │
        │  consumers (never import generated/ directly)
        ▼
ibmverify-cli                ← imports client/, uses c.Apps / c.Users / c.APIClients
terraform-provider-verify    ← imports client/
```

**The rule:** never import `generated/` directly from the CLI or Terraform provider. Always go through `client.Client` which exposes `c.Apps`, `c.Users`, and `c.APIClients`.

---

## Quick start

```go
import "github.com/ChrisVerde02/ibmverify-go/client"

c, err := client.New("https://example.verify.ibm.com",
    client.WithClientCredentials("your-client-id", "your-client-secret"),
)
```

One `client.Client` gives you access to all domains:

```go
// Token operations
token,  err := c.Token.ClientCredentials(ctx)
token,  err := c.Token.Exchange(ctx, signedJWT)
info,   err := c.Token.Introspect(ctx, accessToken)

// Signer certificates
cert,   err := c.Certs.Get(ctx, "demotokensigner")
        err  = c.Certs.Import(ctx, "demotokensigner", pemString)
        err  = c.Certs.Delete(ctx, "demotokensigner")

// Applications
apps,   err := c.Apps.List(ctx, nil)
app,    err := c.Apps.Get(ctx, applicationID)
result, err := c.Apps.Create(ctx, &generated.ApplicationRequestBean{...})
        err  = c.Apps.Delete(ctx, applicationID)

// Users (SCIM v2)
users,  err := c.Users.List(ctx, nil)                                    // all users
users,  err := c.Users.List(ctx, &generated.GetUsersRequest{Filter: &f}) // SCIM filter
user,   err := c.Users.Get(ctx, userID)

// API Clients (Dynamic Client Registration)
clients, err := c.APIClients.List(ctx, nil)
client,  err := c.APIClients.Get(ctx, clientID)
result,  err := c.APIClients.Create(ctx, &generated.APIClientConfigRequest{...})
         err  = c.APIClients.Delete(ctx, clientID)
```

---

## Packages

### `client` — top-level client and authentication

```go
import "github.com/ChrisVerde02/ibmverify-go/client"

c, err := client.New(tenantURL, client.WithClientCredentials(id, secret))
```

**Options:**

| Option | Description |
|---|---|
| `WithClientCredentials(id, secret)` | OAuth client ID and secret |
| `WithHTTPClient(hc)` | Replace the default HTTP client |
| `WithTimeout(d)` | Set request timeout (default 30s) |

**Token operations (`c.Token`):**

| Method | Description |
|---|---|
| `ClientCredentials(ctx)` | Acquire an access token via client credentials grant |
| `Exchange(ctx, jwt)` | RFC 8693 token exchange — swap a signed JWT for an IBM Verify token |
| `Introspect(ctx, token)` | Inspect an access token — returns subject, username, scope, expiry |

**Certificate operations (`c.Certs`):**

| Method | Description |
|---|---|
| `Import(ctx, label, pem)` | Upload a signer certificate |
| `Get(ctx, label)` | Fetch a signer certificate by label |
| `Delete(ctx, label)` | Delete a signer certificate by label |

**Error handling:**

```go
var apiErr *client.APIError
if errors.As(err, &apiErr) {
    apiErr.IsNotFound()   // HTTP 404
    apiErr.IsAuth()       // HTTP 401/403
    apiErr.IsRateLimit()  // HTTP 429
    apiErr.IsRetryable()  // 429 or 5xx — automatically retried
}

if errors.Is(err, client.ErrNotFound) { ... }
```

The SDK automatically retries `429` and `5xx` responses up to 3 times with exponential backoff (1s, 2s).

---

### `apps` — application management

```go
// accessed via c.Apps — do not instantiate directly
apps,   err := c.Apps.List(ctx, nil)
app,    err := c.Apps.Get(ctx, "application-id")
result, err := c.Apps.Create(ctx, &generated.ApplicationRequestBean{
    Name:       "My App",
    TemplateID: "template-id",
})
err = c.Apps.Delete(ctx, "application-id")
```

`List` and `Get` return `[]map[string]interface{}` and `map[string]interface{}` respectively. IBM Verify returns `applicationState` as a `bool` in live responses despite the spec saying `string` — the raw HTTP approach handles this transparently.

---

### `users` — user management (SCIM v2)

```go
// accessed via c.Users — do not instantiate directly

// list all users
users, err := c.Users.List(ctx, nil)

// list with SCIM filter
filter := `userName eq "john"`
users, err := c.Users.List(ctx, &generated.GetUsersRequest{Filter: &filter})

// get by ID
user, err := c.Users.Get(ctx, "user-id")
```

`List` and `Get` return raw JSON maps. IBM Verify requires `Accept: application/scim+json` — this is set automatically.

---

### `apiclients` — dynamic client registration (DCR)

```go
// accessed via c.APIClients — do not instantiate directly

// list all API clients
clients, err := c.APIClients.List(ctx, nil)

// get one by client ID
cl, err := c.APIClients.Get(ctx, "client-id")

// create a new API client
result, err := c.APIClients.Create(ctx, &generated.APIClientConfigRequest{
    ClientName:   "My Automation Client",
    Entitlements: []string{"manageOAuthClients"},
    Enabled:      true,
})
// result["clientSecret"] is only present here — IBM never returns it again

// delete
err = c.APIClients.Delete(ctx, "client-id")
```

All methods return `map[string]interface{}`. Array fields (`entitlements`, `ipFilters`) are `[]interface{}` after the JSON round-trip — assert accordingly:

```go
if ents, ok := result["entitlements"].([]interface{}); ok {
    for _, e := range ents {
        fmt.Println(e.(string))
    }
}
```

> **`clientSecret` is only returned by `Create`.** `Get` and `List` never include it. Capture it at creation time.

---

### `crypto` — local cryptography

```go
import "github.com/ChrisVerde02/ibmverify-go/crypto"

// Generate a self-signed certificate
cert, err := crypto.GenerateSelfSignedCertificate(crypto.CertificateRequest{
    CommonName:   "DemoTokenSigner",
    Organization: "IBM",
    Country:      "US",
    ValidityDays: 365,
    KeySize:      4096,
})
// cert.CertificatePEM, cert.PrivateKeyPEM

// Sign a JWT
jwt, err := crypto.GenerateSignedJWT(crypto.JWTRequest{
    Issuer:        "https://example.ibm.com",
    Subject:       "user@example.com",
    KeyID:         "demotokensigner",
    PrivateKeyPEM: cert.PrivateKeyPEM,
    ExpiresIn:     15 * time.Minute,
})
// jwt.Token
```

---

## Full flow example

```go
c, err := client.New("https://example.verify.ibm.com",
    client.WithClientCredentials("client-id", "client-secret"),
)

// 1. Get the signer cert
cert, err := c.Certs.Get(ctx, "demotokensigner")

// 2. Sign a JWT locally
jwt, err := crypto.GenerateSignedJWT(crypto.JWTRequest{
    Issuer:        "https://example.ibm.com",
    Subject:       "user@example.com",
    KeyID:         "demotokensigner",
    PrivateKeyPEM: cert.PrivateKeyPEM,
    ExpiresIn:     5 * time.Minute,
})

// 3. Exchange for an IBM Verify access token
token, err := c.Token.Exchange(ctx, jwt.Token)
fmt.Println(token.AccessToken)

// 4. List applications
apps, err := c.Apps.List(ctx, nil)
for _, app := range apps {
    fmt.Printf("%s  %s\n", app["id"], app["name"])
}
```

---

## Spec regeneration

When IBM updates their API spec:

```bash
# 1. Sanitize the raw spec
./specs/sanitize.sh

# 2. Regenerate (requires Docker and Fern CLI)
fern generate --group go-sdk

# 3. Fix Go version if Fern bumped it
sed -i '' 's/^go 1\.26\..*/go 1.23/' go.mod

# 4. Verify
go build ./...
go test ./...
```

> `generated/` tests that require WireMock (`generated/*/..._test/`) will fail without a running WireMock container — this is expected. Only the tests outside `generated/` need to pass.

---

## Development

```bash
# Run all non-WireMock tests
go test $(go list ./... | grep -v 'generated/.*/.*_test')

# Run everything (WireMock tests will fail without Docker — that's expected)
go test ./...

# Build
go build ./...

# Vet
go vet ./...
```

## Requirements

- Go 1.23+
- [Fern CLI](https://buildwithfern.com) + Docker — only needed to regenerate `generated/`

## Versioning

Follows [semantic versioning](https://semver.org). Current version: **v1.6.4**

| Change | Version bump |
|---|---|
| Breaking change to `client/`, `apps/`, `users/`, `apiclients/`, or `crypto/` public API | Major |
| New method or domain package | Minor |
| Bug fix, test addition, doc update | Patch |
