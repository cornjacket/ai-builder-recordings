# lifecycle

Purpose: Provides in-memory user CRUD and token-based authentication for the IAM listener. Implements five HTTP endpoints exposed via an `*http.ServeMux` for registration by integrate-iam.

Tags: architecture, design

## File Index

| File | Description |
|------|-------------|
| `lifecycle.go` | `UserStore`, `TokenStore`, `handler` struct, `Handler()` constructor, all five HTTP handler methods |
| `lifecycle_test.go` | Table-driven HTTP tests for all five endpoints using `net/http/httptest` |

## Overview

The package owns two in-memory stores, each protected by a `sync.RWMutex`:

- **UserStore** — two maps keyed by ID and by username, holding `{id, username, passwordHash}` records.
- **TokenStore** — one map keyed by opaque token string, mapping to the owning user ID.

IDs and tokens are generated with `github.com/google/uuid` (v4). Passwords are hashed with `golang.org/x/crypto/bcrypt` (cost 12). All response bodies are JSON; errors use `{"error":"<message>"}`.

`Handler()` is the single public entry point. It constructs both stores, wires them into a private `handler` struct, and returns a `*http.ServeMux` with all five routes registered. Because the module targets Go 1.21, method-prefix patterns (`POST /users`) are not available — the mux registers bare paths and each handler switches on `r.Method`. Path parameters for `/users/{id}` are extracted by splitting `r.URL.Path`.

### Dependencies

- `github.com/google/uuid` — already in `go.mod`; used for ID and token generation
- `golang.org/x/crypto/bcrypt` — must be added to `go.mod` (`go get golang.org/x/crypto`)

See [api.md](api.md) for the full endpoint contract.

## Documentation

### Design
| File | Description |
|------|-------------|
| api.md | Defines the HTTP request and response contract for all five lifecycle endpoints |

