# Contributing — Adding a New IBM Verify Domain

This guide walks through adding a completely new IBM Verify domain to the SDK
(e.g. Groups, Agents, Workflows, Policies). After following these steps the new
domain will be available as `c.YourDomain` on the top-level `client.Client`,
with full CLI and Terraform support wired in.

The same pattern was used for `apps/`, `users/`, and `apiclients/`. Every step
below references those packages as the working example.

---

## The four-layer architecture

Before touching any file, understand what layer does what:

```
specs/openapi_trimmed.yaml   ← you add your new paths here (Step 1)
        │
        │  fern generate
        ▼
generated/                   ← Fern output — DO NOT EDIT by hand
  types.go                   ← request/response structs appear here automatically
  yourDomain/client.go       ← generated typed client appears here automatically
        │
        │  thin handwritten wrapper (Step 3)
        ▼
yourdomain/yourdomain.go     ← absorbs IBM quirks; what CLI + Terraform import
        │
        │  wired onto top-level Client (Step 4)
        ▼
client/client.go             ← c.YourDomain exposed here
        │
        ▼
ibmverify-cli                ← cmd/yourdomain/ (Step 5)
terraform-provider-verify    ← internal/resources/yourdomain_resource.go (Step 6)
```

**The rule:** the CLI and Terraform provider never import `generated/` directly.
They always go through `client.Client` → `c.YourDomain`.

---

## Step 1 — Add your paths to the trimmed spec

Open `specs/openapi_trimmed.yaml`. This is the 6-path file Fern reads.
Add your new paths under the existing `paths:` key.

**Why not the full spec?** Fern v1.57.2 silently truncates `types.go` when fed
IBM's full 385-path spec. The trimmed file is the workaround. The full spec is
preserved in `specs/openapi_corrected.yaml` — copy paths from there.

### What to copy

Find your domain in `specs/openapi_corrected.yaml`:

```bash
grep -n "your-domain-path" specs/openapi_corrected.yaml
```

Copy the relevant `paths:` entries into `specs/openapi_trimmed.yaml`.
Only include the HTTP verbs you need (GET, POST, DELETE — skip PUT/PATCH
if you don't intend to use them, to keep generation output small).

**Minimal example — adding a read-only Groups domain:**

```yaml
# in specs/openapi_trimmed.yaml, under paths:

  /v2.0/Groups:
    get:
      tags:
        - Groups Management
      summary: Retrieve all groups.
      operationId: getGroups
      parameters:
        - name: filter
          in: query
          schema:
            type: string
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GroupListResponse'

  /v2.0/Groups/{groupId}:
    get:
      tags:
        - Groups Management
      summary: Retrieve a single group.
      operationId: getGroupById
      parameters:
        - name: groupId
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Group'
```

Also copy any `$ref` schemas those paths use into the `components/schemas:`
section of the trimmed spec. Only the schemas your paths reference are needed.

---

## Step 2 — Add overlay entries

Open `fern/verify-overlay.yml`. Add a section for your domain.

IBM's `operationId` values are verbose (`GetGroupsV2`, `GetGroupByIdV2`).
The overlay renames them to clean Go method names.

```yaml
# in fern/verify-overlay.yml

  # ── Groups ──────────────────────────────────────────────────────────────────
  - target: "$.paths['/v2.0/Groups'].get"
    update:
      x-fern-sdk-group-name: groups
      x-fern-sdk-method-name: list

  - target: "$.paths['/v2.0/Groups/{groupId}'].get"
    update:
      x-fern-sdk-group-name: groups
      x-fern-sdk-method-name: get
```

`x-fern-sdk-group-name` becomes the Go package name under `generated/`
(e.g. `generated/groups/`). `x-fern-sdk-method-name` becomes the Go method
name on that client.

---

## Step 3 — Regenerate

```bash
# Regenerate generated/ from the updated trimmed spec
fern generate --group go-sdk

# Fix Go version if Fern bumped it
sed -i '' 's/^go 1\.26\..*/go 1.23/' go.mod

# Verify it compiles
go build ./...
```

**Requires:** Docker running, Fern CLI installed (`npm install -g fern-api`).

After generation you will see:
- `generated/groups/client.go` — typed `List(ctx, ...)` and `Get(ctx, ...)` methods
- Any new types added to `generated/types.go`

> Tests in `generated/groups/groups_test/` will fail without a WireMock container
> on `localhost:8080`. This is expected — ignore them permanently.

---

## Step 4 — Write the wrapper package

Create `groups/groups.go`. Copy the structure from `apps/apps.go` and adapt it.

The wrapper's job is to absorb IBM quirks so the CLI and Terraform never see them.

### Decide: typed client or raw HTTP?

| Use the **Fern generated client** | Use **raw HTTP** |
|---|---|
| IBM's live response matches the spec | IBM's live response has type mismatches |
| Response body has proper JSON schemas | Endpoint has no response schema in spec |
| Simple request/response shapes | Response needs custom unwrapping |

Check by calling the real endpoint with `--debug` before deciding:

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  "https://YOUR-TENANT.verify.ibm.com/v2.0/Groups" | python3 -m json.tool | head -30
```

Compare the live field types to what Fern generated in `generated/types.go`.
If they match → typed client. If they don't → raw HTTP.

**Known IBM quirks by domain:**

| Domain | Quirk | How it's handled |
|---|---|---|
| `apps` | `applicationState` is `bool` in live response, `string` in spec | `List`/`Get` use raw HTTP |
| `users` | Requires `Accept: application/scim+json` (returns 406 without it) | Set header on every request |
| `apiclients` | Fern's raw client returns `Body: nil` for DCR endpoints | All ops use raw HTTP |

### Wrapper skeleton

```go
// Package groups provides a high-level client for IBM Verify group management.
package groups

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// Client manages IBM Verify groups.
type Client struct {
    tenantURL  string
    getToken   func(ctx context.Context) (string, error)
    httpClient *http.Client
}

// New returns a Groups Client.
func New(tenantURL string, getToken func(ctx context.Context) (string, error)) *Client {
    return &Client{
        tenantURL:  tenantURL,
        getToken:   getToken,
        httpClient: &http.Client{},
    }
}

// rawGet performs an authenticated GET and returns the body.
func (c *Client) rawGet(ctx context.Context, path string) ([]byte, error) {
    token, err := c.getToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("groups: get token: %w", err)
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.tenantURL+path, nil)
    if err != nil {
        return nil, fmt.Errorf("groups: create request: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("groups: send request: %w", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
    if err != nil {
        return nil, fmt.Errorf("groups: read response: %w", err)
    }
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("groups: HTTP %d: %s", resp.StatusCode, string(body))
    }
    return body, nil
}

// List returns all groups as raw JSON maps.
func (c *Client) List(ctx context.Context) ([]map[string]interface{}, error) {
    body, err := c.rawGet(ctx, "/v2.0/Groups")
    if err != nil {
        return nil, fmt.Errorf("groups: list: %w", err)
    }
    // Unwrap IBM's envelope — check the actual live response shape first
    var wrapper struct {
        Resources []map[string]interface{} `json:"Resources"`
    }
    if json.Unmarshal(body, &wrapper) == nil && wrapper.Resources != nil {
        return wrapper.Resources, nil
    }
    var list []map[string]interface{}
    _ = json.Unmarshal(body, &list)
    return list, nil
}

// Get returns a single group by ID.
func (c *Client) Get(ctx context.Context, groupID string) (map[string]interface{}, error) {
    body, err := c.rawGet(ctx, "/v2.0/Groups/"+groupID)
    if err != nil {
        return nil, fmt.Errorf("groups: get %s: %w", groupID, err)
    }
    var m map[string]interface{}
    if err := json.Unmarshal(body, &m); err != nil {
        return nil, fmt.Errorf("groups: get %s: unmarshal: %w", groupID, err)
    }
    return m, nil
}
```

---

## Step 5 — Wire it onto `client.Client`

Open `client/client.go`. Add three things:

**1. Import:**
```go
import (
    // existing imports ...
    "github.com/ChrisVerde02/ibmverify-go/groups"
)
```

**2. Field on the struct:**
```go
type Client struct {
    // existing fields ...

    // Groups provides methods for IBM Verify group management.
    Groups *groups.Client
}
```

**3. Instantiate in `New()`**, after the existing domain clients:
```go
c.Groups = groups.New(c.tenantURL, getToken)
```

Run `go build ./...` to confirm it compiles.

---

## Step 6 — Write tests

Create `groups/groups_test.go`. Use `net/http/httptest` — no live tenant needed.

Pattern from `apiclients/apiclients_test.go`:

```go
func TestGroupsList(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/v1.0/endpoint/default/token":
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprint(w, `{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`)
        case "/v2.0/Groups":
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprint(w, `{"Resources":[{"id":"g1","displayName":"Admins"}]}`)
        default:
            http.NotFound(w, r)
        }
    }))
    defer srv.Close()

    c := groups.New(srv.URL, func(ctx context.Context) (string, error) {
        return "test-token", nil
    })

    list, err := c.List(context.Background())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(list) != 1 {
        t.Fatalf("expected 1 group, got %d", len(list))
    }
    if list[0]["displayName"] != "Admins" {
        t.Errorf("expected Admins, got %v", list[0]["displayName"])
    }
}
```

Cover: success, 404, 401, empty list.

Run: `go test ./groups/...`

---

## Step 7 — Add the CLI command (ibmverify-cli)

In `ibmverify-cli`, create `cmd/group/`. Copy `cmd/app/` as the starting point
and change `app` → `group` throughout.

Each file follows the same pattern:

```
cmd/group/
  group.go       ← parent Cobra command (GroupCmd)
  list.go        ← ibmverify group list
  get.go         ← ibmverify group get --id <id>
  group_test.go  ← httptest-backed tests for all subcommands
```

Register it in `cmd/ibmverify/main.go`:

```go
import "github.com/ChrisVerde02/ibmverify-cli/cmd/group"

func init() {
    rootCmd.AddCommand(group.GroupCmd)
}
```

Run `go test ./cmd/group/...` before committing.

---

## Step 8 — Add the Terraform resource (terraform-provider-verify)

In `terraform-provider-verify`, create
`internal/resources/group_resource.go`. Copy
`internal/resources/application_resource.go` as the starting point.

Key things to adapt:

1. **Schema** — define your attributes in `Schema()`. Mark computed fields as
   `Computed: true`, required inputs as `Required: true`.

2. **Idempotency** — in `Create()`, check if a resource with the same key
   already exists before creating. If it does, adopt it into state instead of
   returning an error. Pattern:

   ```go
   // in Create — check for existing before POST
   existing, err := r.client.Groups.List(ctx)
   for _, g := range existing {
       if g["displayName"] == plan.Name.ValueString() {
           // adopt — write ID to state and return
           state.ID = types.StringValue(g["id"].(string))
           resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
           return
       }
   }
   ```

3. **Import** — implement `ImportState()`:

   ```go
   func (r *GroupResource) ImportState(ctx context.Context,
       req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
       resource.ImportStatePassthroughID(ctx, path.Root("group_id"), req, resp)
   }
   ```

4. **Error handling** — use `client.APIError` consistently:

   ```go
   var apiErr *client.APIError
   if errors.As(err, &apiErr) && apiErr.IsNotFound() {
       resp.State.RemoveResource(ctx)
       return
   }
   resp.Diagnostics.AddError("Failed to read group", err.Error())
   ```

5. **Register** in `internal/provider/provider.go`:

   ```go
   func (p *VerifyProvider) Resources(ctx context.Context) []func() resource.Resource {
       return []func() resource.Resource{
           // existing resources ...
           resources.NewGroupResource,
       }
   }
   ```

Run `go test ./internal/resources/tests/...` before committing.

---

## Checklist

Before opening a PR, confirm all of these:

- [ ] New paths added to `specs/openapi_trimmed.yaml`
- [ ] Overlay entries added to `fern/verify-overlay.yml`
- [ ] `fern generate` run, `go.mod` Go version fixed if needed
- [ ] Wrapper package created under `yourdomain/yourdomain.go`
- [ ] Domain wired onto `client.Client` — struct field + `New()` instantiation
- [ ] `go build ./...` passes
- [ ] `go test $(go list ./... | grep -v 'generated/.*/.*_test')` all pass
- [ ] CLI command added to `ibmverify-cli/cmd/yourdomain/`
- [ ] Terraform resource added to `terraform-provider-verify/internal/resources/`
- [ ] Idempotency implemented in Terraform Create (adopt-on-duplicate)
- [ ] `terraform import` implemented
- [ ] README updated in all three repos

---

## Common problems

### Fern generates nothing for my new paths

Check that your `operationId` values in the trimmed spec are unique and that
the `$ref` schemas exist in `components/schemas:`. Fern silently skips
operations with missing schema references.

### Generated struct field types don't match IBM's live response

This is the most common issue. IBM's live API and their OpenAPI spec diverge.
Use raw HTTP for that operation — see `apps/apps.go` `List()` for the pattern.
Put a comment explaining the mismatch so the next person knows why.

### IBM returns 406 on my new endpoint

IBM requires `Accept: application/scim+json` on SCIM endpoints (`/v2.0/`).
Set it explicitly on every request in your wrapper:

```go
req.Header.Set("Accept", "application/scim+json")
```

### `go.mod` says `go 1.26.x` after fern generate

Fern bumps the Go version. Fix it:

```bash
sed -i '' 's/^go 1\.26\..*/go 1.23/' go.mod
```

### WireMock tests fail

Expected. `generated/yourdomain/yourdomain_test/` tests require a live
WireMock container on `localhost:8080`. Ignore them permanently. Only the
tests in `yourdomain/yourdomain_test.go` (your handwritten httptest tests)
need to pass.
