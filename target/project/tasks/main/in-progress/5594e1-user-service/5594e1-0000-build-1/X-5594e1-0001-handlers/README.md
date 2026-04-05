# handlers

## Goal

HTTP CRUD handlers for the user management API. Accepts a store interface and returns an http.Handler (or registers on a mux). Full API contract — POST /users {"name":string,"email":string} → 201 {"id":string,"name":string,"email":string}; GET /users/{id} → 200 {"id":string,"name":string,"email":string} or 404 {}; PUT /users/{id} {"name":string,"email":string} → 200 {"id":string,"name":string,"email":string} or 404 {}; DELETE /users/{id} → 204 (no body) or 404 {}. All request and response bodies are JSON. Returns 400 on malformed JSON request bodies.

## Context

### Level 1 — 5594e1-0000-build-1


## Design

**Language:** Go, standard library only (`net/http`, `encoding/json`).

**Package:** `package handlers` at `internal/userservice/handlers`.

**Module:** `github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service`

**Dependencies:** Only `internal/userservice/store` (same module) for the `store.User` type.

### Files

- `handlers.go` — all production code
- `handlers_test.go` — tests using `net/http/httptest`

### handlers.go structure

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

// Store is the interface the handlers depend on.
type Store interface {
    Create(name, email string) store.User
    Get(id string) (store.User, bool)
    Update(id, name, email string) (store.User, bool)
    Delete(id string) bool
}

// New registers all routes and returns the mux.
func New(s Store) http.Handler { ... }

// unexported request/response types
type createRequest  struct { Name string `json:"name"`; Email string `json:"email"` }
type updateRequest  struct { Name string `json:"name"`; Email string `json:"email"` }
type userResponse   struct { ID string `json:"id"`; Name string `json:"name"`; Email string `json:"email"` }
```

### Route registration (Go 1.22+ ServeMux patterns)

```
mux.HandleFunc("POST /users",        handleCreate(s))
mux.HandleFunc("GET /users/{id}",    handleGet(s))
mux.HandleFunc("PUT /users/{id}",    handleUpdate(s))
mux.HandleFunc("DELETE /users/{id}", handleDelete(s))
```

### Handler behaviour

**handleCreate:**
1. `json.NewDecoder(r.Body).Decode(&req)` — on error → `w.WriteHeader(400)`; return
2. `s.Create(req.Name, req.Email)` → user
3. `w.Header().Set("Content-Type", "application/json")`, `w.WriteHeader(201)`
4. `json.NewEncoder(w).Encode(userResponse{...})`

**handleGet:**
1. `id := r.PathValue("id")`
2. `u, ok := s.Get(id)` — if !ok → write 404 + `{}`; return
3. Write 200 + userResponse JSON

**handleUpdate:**
1. `id := r.PathValue("id")`
2. `json.NewDecoder(r.Body).Decode(&req)` — on error → write 400; return
3. `u, ok := s.Update(id, req.Name, req.Email)` — if !ok → write 404 + `{}`; return
4. Write 200 + userResponse JSON

**handleDelete:**
1. `id := r.PathValue("id")`
2. `ok := s.Delete(id)` — if !ok → write 404 + `{}`; return
3. `w.WriteHeader(204)` (no body)

### 404 helper

```go
func writeNotFound(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("{}"))
}
```

### JSON encoding note

Use `json.NewEncoder(w).Encode(v)` for responses (appends newline, which is acceptable). Set `Content-Type: application/json` before calling `WriteHeader`.

### Constraints

- Go 1.22+ required for `r.PathValue` and method-qualified mux patterns.
- No third-party dependencies.
- The `Store` interface is defined in this package — the caller passes a `*store.Store` which satisfies it structurally.

## Acceptance Criteria

1. `POST /users` with body `{"name":"Alice","email":"alice@example.com"}` returns HTTP 201 and a JSON body containing `"name":"Alice"`, `"email":"alice@example.com"`, and a non-empty `"id"` field.
2. `POST /users` with body `{invalid` returns HTTP 400.
3. `GET /users/{id}` where `{id}` is the ID from criterion 1 returns HTTP 200 and the same JSON body (`id`, `name`, `email` matching).
4. `GET /users/{id}` where `{id}` is an unknown UUID returns HTTP 404 with body `{}`.
5. `PUT /users/{id}` (existing ID) with body `{"name":"Bob","email":"bob@example.com"}` returns HTTP 200 and JSON body with `"name":"Bob"`, `"email":"bob@example.com"`, and the same `id`.
6. `PUT /users/{id}` (unknown ID) with valid body returns HTTP 404 with body `{}`.
7. `PUT /users/{id}` with body `{invalid` returns HTTP 400.
8. `DELETE /users/{id}` (existing ID) returns HTTP 204 with an empty body.
9. `DELETE /users/{id}` (unknown ID) returns HTTP 404 with body `{}`.
10. All non-204 JSON responses carry `Content-Type: application/json`.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/output && go test ./internal/userservice/handlers/...
```

