# integrate-user-service

## Goal

Wire all components into a cohesive unit and verify this level's acceptance criteria. Write main.go at the output root: instantiate the store, pass it to the handlers, register routes on net/http.ServeMux (POST /users, GET /users/{id}, PUT /users/{id}, DELETE /users/{id}), and call http.ListenAndServe(":8080", mux). End-to-end acceptance criteria: POST /users returns 201 with generated id; GET /users/{id} returns 200 with correct record or 404; PUT /users/{id} returns 200 with updated record or 404; DELETE /users/{id} returns 204 or 404; server listens on port 8080; all responses are JSON.

## Context

### Level 1 — 5594e1-0000-build-1


## Design

**Language:** Go (module `github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service`)

**Files to produce** (both at the output root `/Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/output`):

### `main.go`

```go
package main

import (
    "net/http"

    "github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/handlers"
    "github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

func main() {
    s := store.New()
    h := handlers.New(s)
    http.ListenAndServe(":8080", h)
}
```

`*store.Store` satisfies `handlers.Store` directly — no adapter needed.

### `main_test.go`

`package main`, integration test using `httptest.NewServer(handlers.New(store.New()))`. Tests never bind port 8080; they hit the test server's URL. The test covers:

1. POST /users → 201, response body has `id`, `name`, `email`
2. GET /users/{id} (known) → 200, correct fields
3. PUT /users/{id} (known) → 200, updated fields
4. GET /users/{id} (verify update stuck) → 200
5. DELETE /users/{id} → 204
6. GET /users/{id} (after delete) → 404
7. DELETE /users/{id} (after delete) → 404
8. GET /users/nonexistent → 404
9. PUT /users/nonexistent → 404

No external test libraries; use `encoding/json`, `net/http`, `strings`, `testing`.

**Dependencies:** stdlib only + two internal packages already implemented.

**Constraints:**
- `handlers.New` already registers all routes; `main.go` must not re-register them.
- The `*store.Store` concrete type satisfies the `handlers.Store` interface without any casting — pass it directly.
- `http.ListenAndServe` return value can be ignored or logged; silence is fine for this service.

## Acceptance Criteria

1. `go build ./...` succeeds with no errors.
2. POST /users with `{"name":"Alice","email":"alice@example.com"}` returns HTTP 201 with a JSON body containing a non-empty `id` string, `"name":"Alice"`, and `"email":"alice@example.com"`.
3. GET /users/{id} using the id from criterion 2 returns HTTP 200 with the same `id`, `name`, and `email`.
4. GET /users/{id} with an id that was never created returns HTTP 404.
5. PUT /users/{id} using the id from criterion 2 with `{"name":"Bob","email":"bob@example.com"}` returns HTTP 200 with `"name":"Bob"` and `"email":"bob@example.com"`.
6. PUT /users/{id} with an id that was never created returns HTTP 404.
7. DELETE /users/{id} using the id from criterion 2 returns HTTP 204.
8. DELETE /users/{id} using the same id a second time returns HTTP 404.
9. GET /users/{id} using the deleted id returns HTTP 404.
10. All non-204 responses have `Content-Type: application/json`.
11. `go test ./...` passes with all tests green.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/output && go test ./...
```

