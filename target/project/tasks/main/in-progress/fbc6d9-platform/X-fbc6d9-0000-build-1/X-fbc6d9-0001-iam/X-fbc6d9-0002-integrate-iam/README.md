# integrate-iam

## Goal

Wire lifecycle and authz handlers into a single http.ServeMux and expose it for the platform main.go. Create iam.go (package iam) in the output directory. It must instantiate lifecycle.New() and authz.New(), register their routes on a shared mux (route prefixes: /users and /auth go to lifecycle; /roles, /authz, and /users/ sub-paths for roles go to authz — use pattern matching carefully so /users/{id}/roles routes to authz while /users/{id} routes to lifecycle), and return the mux via a NewMux() function. No port binding here — the platform main.go calls http.ListenAndServe(":8082", iam.NewMux()). Task Level: INTERNAL.

## Context

### Level 1 — fbc6d9-0000-build-1

### Level 2 — fbc6d9-0001-iam
Identity and access management listener (port 8082) composed of two sub-components — lifecycle (user CRUD and authentication) and authz-rbac (roles and permission checks) — that together handle all ten IAM HTTP endpoints.

## Design

**Language:** Go (module `github.com/cornjacket/platform-monolith`, `go 1.25.0`)

**File to produce:** `iam.go` in the output directory (`internal/iam/iam.go` relative to module root)

**Package declaration:** `package iam`

**Dependencies:**
- `net/http` — stdlib
- `strings` — stdlib (for `strings.HasSuffix`)
- `github.com/cornjacket/platform-monolith/internal/iam/internal/iam/lifecycle` — exposes `Handler() *http.ServeMux`
- `github.com/cornjacket/platform-monolith/internal/iam/internal/iam/authz` — exposes `Handler() *http.ServeMux`

**Public API:**

```go
func NewMux() *http.ServeMux
```

**Implementation:**

```go
func NewMux() *http.ServeMux {
    lc := lifecycle.Handler()
    az := authz.Handler()

    mux := http.NewServeMux()

    // lifecycle routes
    mux.Handle("/users", lc)
    mux.Handle("/auth/login", lc)
    mux.Handle("/auth/logout", lc)

    // authz routes
    mux.Handle("/roles", az)
    mux.Handle("/authz/check", az)

    // /users/{id} → lifecycle; /users/{id}/roles → authz
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/roles") {
            az.ServeHTTP(w, r)
        } else {
            lc.ServeHTTP(w, r)
        }
    })

    return mux
}
```

**Constraints:**
- No `init()`, no global state, no port binding.
- Both sub-muxes receive the original unmodified request; their internal path-extraction (TrimPrefix, Split) works on the full path.
- The `/users/` suffix check (`/roles`) is safe because no other valid path under `/users/` ends in `/roles` except the role-assignment sub-resource.
- Go `internal` visibility: `iam.go` lives at `internal/iam/`, which is the tree root for `internal/iam/internal/`, so both imports are permitted by the compiler.

**Test file to produce:** `iam_test.go` in the same directory. Use `net/http/httptest` to drive `NewMux()`. Each test:
1. Exercises the full happy-path round-trip (e.g. create user, then get user) to confirm routing reaches the correct handler.
2. Confirms `/users/{id}/roles` reaches authz (not lifecycle).
3. Confirms `/users/{id}` reaches lifecycle (not authz).

## Acceptance Criteria

1. `NewMux()` returns a non-nil `*http.ServeMux` with no panics.
2. `POST /users` with `{"username":"u1","password":"pw"}` → 201 with `{"id":…,"username":"u1"}`.
3. `GET /users/{id}` with a valid id returned from criterion 2 → 200 with `{"id":…,"username":"u1"}`.
4. `DELETE /users/{id}` with a valid id → 200 (lifecycle, not authz).
5. `POST /auth/login` with valid credentials → 200 with `{"token":…}`.
6. `POST /auth/logout` with valid `Authorization: Bearer <token>` → 200.
7. `POST /roles` with `{"name":"admin","permissions":["read"]}` → 201 with `{"id":…,"name":"admin","permissions":["read"]}`.
8. `GET /roles` → 200 with a JSON array (authz handler reached, not lifecycle).
9. `POST /users/{id}/roles` with a valid roleId → 201 (authz handler reached, not lifecycle).
10. `GET /users/{id}/roles` → 200 with a JSON array (authz handler reached).
11. `POST /authz/check` with `{"userId":…,"permission":"read"}` → 200 with `{"allowed":true}` after assigning a role that has `"read"` permission.
12. `GET /users/{id}` does NOT return a response handled by authz — specifically, it returns `{"id":…,"username":…}` (lifecycle shape), not a role array.
13. `GET /users/{id}/roles` does NOT return a response handled by lifecycle — it returns a role array, not a user object.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/platform-monolith/output && go test ./internal/iam/...
```

