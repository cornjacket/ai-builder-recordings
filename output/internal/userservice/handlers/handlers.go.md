# handlers.go

Purpose: Implements HTTP CRUD handlers for the user service, wiring POST/GET/PUT/DELETE `/users` routes onto a `*http.ServeMux` via a `Store` interface. JSON request bodies are decoded on ingress and JSON responses are encoded on egress.

Tags: api, interface

## Public API

### `Store` interface

```go
type Store interface {
    Create(user store.User) store.User
    Get(id string) (store.User, bool)
    Update(id string, user store.User) (store.User, bool)
    Delete(id string) bool
}
```

Declares the four data-access operations that a backing store must satisfy. The `store` package's `*store.Store` implements this interface, but any conforming type may be substituted (useful for testing).

### `Handler` struct

```go
type Handler struct { ... }
```

Holds a `Store` reference and exposes the route-registration method. Constructed via `New`.

### `New`

```go
func New(s Store) *Handler
```

Returns a `*Handler` backed by `s`. `s` must not be nil.

### `RegisterRoutes`

```go
func (h *Handler) RegisterRoutes(mux *http.ServeMux)
```

Registers four routes onto `mux` using Go 1.22 method-prefixed patterns:

| Method | Pattern | Handler |
|--------|---------|---------|
| POST | `/users` | `createUser` |
| GET | `/users/{id}` | `getUser` |
| PUT | `/users/{id}` | `updateUser` |
| DELETE | `/users/{id}` | `deleteUser` |

## Key Internals

### Request/response shapes

`userInput` (Name, Email) is decoded from POST and PUT request bodies; the caller never supplies an ID. `userJSON` (ID, Name, Email) is the uniform response envelope for all successful write and read operations.

### `createUser` / `getUser` / `updateUser` / `deleteUser`

Each handler follows the same pattern: decode the path value or request body, delegate to the store, and write a typed JSON response or a 404/400 error. `deleteUser` returns `204 No Content` on success with no body.

### `writeJSON` / `writeError`

`writeJSON` sets `Content-Type: application/json`, writes the status code, then streams `v` via `json.NewEncoder`. `writeError` wraps `writeJSON` with a `{"error": "..."}` envelope, keeping error responses structurally consistent.

## Dependencies

| Package | Role |
|---------|------|
| `encoding/json` | JSON decode/encode for request bodies and responses |
| `net/http` | HTTP handler types, path value extraction, status constants |
| `github.com/cornjacket/ai-builder/.../store` | `store.User` type consumed by the `Store` interface |
