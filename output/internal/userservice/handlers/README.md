# handlers

Purpose: HTTP CRUD handlers for the user management API, wired to a Store interface.
Handles JSON encoding/decoding, route registration, and HTTP status code mapping for all four user endpoints.

Tags: architecture, design

## Files

| File | Description |
|------|-------------|
| `handlers.go` | `New(Store) http.Handler`, request/response types, and all four route handlers |
| `handlers_test.go` | Table-driven tests for every endpoint covering success and error paths |

## Overview

`New(s Store) http.Handler` constructs a `net/http.ServeMux`, registers four routes using Go 1.22+ enhanced pattern matching (`POST /users`, `GET /users/{id}`, `PUT /users/{id}`, `DELETE /users/{id}`), and returns the mux.

A `Store` interface is defined in this package, matching the concrete `*store.Store` methods:

```go
type Store interface {
    Create(name, email string) store.User
    Get(id string) (store.User, bool)
    Update(id, name, email string) (store.User, bool)
    Delete(id string) bool
}
```

**Request/response types** defined in this package (unexported):

- `createRequest` / `updateRequest`: `Name string json:"name"`, `Email string json:"email"`
- `userResponse`: `ID string json:"id"`, `Name string json:"name"`, `Email string json:"email"`

**Status code mapping:**

| Endpoint | Success | Not Found | Bad JSON |
|----------|---------|-----------|----------|
| POST /users | 201 | — | 400 |
| GET /users/{id} | 200 | 404 + `{}` | — |
| PUT /users/{id} | 200 | 404 + `{}` | 400 |
| DELETE /users/{id} | 204 (no body) | 404 + `{}` | — |

ID extraction uses `r.PathValue("id")` (Go 1.22+ ServeMux). JSON decode errors on request bodies return 400. All JSON responses set `Content-Type: application/json`.

See [api.md](api.md) for the full endpoint contract.
