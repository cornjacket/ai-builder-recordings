# authz

Purpose: Provides in-memory role management and RBAC permission checks for the IAM listener. Implements five HTTP endpoints exposed via an `*http.ServeMux` for registration by integrate-iam.

Tags: architecture, design

## File Index

| File | Description |
|------|-------------|
| `authz.go` | `Role` struct, `Store` (RoleStore + UserRoles), `handler` struct, `Handler()` constructor, all five HTTP handler methods |
| `authz_test.go` | Table-driven HTTP tests for all five endpoints using `net/http/httptest` |

## Overview

The package owns two in-memory stores, both protected by a single `sync.RWMutex`:

- **RoleStore** — `map[string]Role` keyed by role ID, holding `{id, name string, permissions []string}` records.
- **UserRoles** — `map[string][]string` keyed by user ID, holding a slice of role IDs assigned to that user.

Role IDs are generated with `github.com/google/uuid` (v4). All response bodies are JSON; errors use `{"error":"<message>"}`.

`Handler()` is the single public entry point. It constructs the store, wires it into a private `handler` struct, and returns a `*http.ServeMux` with all five routes registered. Because the module targets Go 1.21, method-prefix patterns (`POST /roles`) are not available — the mux registers bare paths and each handler switches on `r.Method`. Path parameters for `/users/{id}/roles` are extracted by splitting `r.URL.Path`.

### Permission Check Logic

`POST /authz/check` looks up the user's assigned role IDs in UserRoles, then for each role ID fetches the Role from RoleStore and scans its `permissions` slice. It returns `{"allowed":true}` on the first match, `{"allowed":false}` if none match or the user has no assigned roles. Roles that no longer exist in RoleStore are silently skipped during the walk.

### 404 Semantics for Role Assignment

`POST /users/{id}/roles` validates that the `roleId` field in the request body refers to an existing role in RoleStore. If the role is not found it returns 404. The `{id}` user path parameter is not validated against any external store — authz has no compile-time dependency on lifecycle.

### Dependencies

- `github.com/google/uuid` — already in `go.mod`; used for role ID generation

See [api.md](api.md) for the full endpoint contract.

## Documentation

### Design
| File | Description |
|------|-------------|
| api.md | Defines the HTTP request and response contract for all five authz endpoints |

### Implementation Notes
| File | Description |
|------|-------------|
| authz.go.md | Documents the design and key decisions of the authz package's role store and RBAC handler |

