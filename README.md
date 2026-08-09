# ibmverify-go

Go SDK for [IBM Verify](https://www.ibm.com/products/verify-identity). Provides typed functions for the IBM Verify OAuth 2.0 APIs — used by [`ibmverify-cli`](https://github.com/ChrisVerde02/ibmverify-cli) and [`terraform-provider-verify`](https://github.com/ChrisVerde02/terraform-provider-verify).

## Packages

### `client` — IBM Verify API calls

| Function | Description |
|---|---|
| `GetClientCredentialsToken` | `POST /v1.0/endpoint/default/token` — client credentials grant |
| `ExchangeToken` | `POST /oauth2/token` — RFC 8693 token exchange |
| `IntrospectToken` | `POST /oauth2/introspect` — token introspection |
| `ImportSignerCert` | `POST /v1.0/signercert` — upload a signer certificate |
| `GetSignerCert` | `GET /v1.0/signercert/{label}` — read a signer certificate |
| `DeleteSignerCert` | `DELETE /v1.0/signercert/{label}` — delete a signer certificate |

### `crypto` — local cryptography

| Function | Description |
|---|---|
| `GenerateSelfSignedCertificate` | Generates an RSA key pair and self-signed X.509 certificate |
| `GenerateSignedJWT` | Signs an RS256 JWT using an RSA private key |

## Installation

```bash
go get github.com/ChrisVerde02/ibmverify-go
```

## Usage

### Client credentials token

```go
import "github.com/ChrisVerde02/ibmverify-go/client"

result, err := client.GetClientCredentialsToken(ctx, client.ClientCredentialsRequest{
    TenantURL:    "https://example.verify.ibm.com",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})
// result.AccessToken, result.ExpiresIn, result.TokenType, result.Scope
```

### Token exchange

```go
result, err := client.ExchangeToken(ctx, client.TokenExchangeRequest{
    TenantURL:        "https://example.verify.ibm.com",
    ClientID:         "your-sts-client-id",
    ClientSecret:     "your-sts-client-secret",
    SubjectToken:     signedJWT,
    SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
})
// result.AccessToken
```

### Token introspection

```go
result, err := client.IntrospectToken(ctx, client.IntrospectionRequest{
    TenantURL:    "https://example.verify.ibm.com",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Token:        accessToken,
})
// result.Active, result.Subject, result.Username, result.ExpiresAt
```

### Generate a self-signed certificate

```go
import "github.com/ChrisVerde02/ibmverify-go/crypto"

cert, err := crypto.GenerateSelfSignedCertificate(crypto.CertificateRequest{
    CommonName:   "DemoTokenSigner",
    Organization: "IBM",
    Country:      "US",
    ValidityDays: 365,
    KeySize:      4096, // 2048, 3072, or 4096
})
// cert.CertificatePEM, cert.PrivateKeyPEM
```

### Sign a JWT

```go
result, err := crypto.GenerateSignedJWT(crypto.JWTRequest{
    Issuer:        "https://demo.ibm.com",
    Subject:       "username",
    KeyID:         "DemoTokenSigner",
    JWTID:         "unique-id",
    PrivateKeyPEM: cert.PrivateKeyPEM,
    ExpiresIn:     15 * time.Minute,
})
// result.Token
```

## Architecture

```
ibmverify-go       ← this repo — all IBM Verify HTTP calls live here
    ↑                   ↑
ibmverify-cli      terraform-provider-verify
(Cobra CLI)        (Terraform provider)
```

The rule: if it makes an HTTP call to IBM Verify, it belongs in this repo. The CLI and provider contain only Cobra/Terraform wiring.

## Requirements

- Go 1.21+

## Versioning

This module follows [semantic versioning](https://semver.org). Current version: **v1.4.0**
