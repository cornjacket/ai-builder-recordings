# integrate-platform

## Goal

Wire the metrics and iam components into a single binary at cmd/platform/main.go. The binary starts the metrics HTTP listener on port 8081 and the IAM HTTP listener on port 8082 in separate goroutines, then blocks until a shutdown signal. There must be exactly one main package in the entire codebase.

End-to-end acceptance tests (net/http against live listeners started in TestMain or a helper) must cover all 12 endpoints verbatim from the spec:

Metrics listener (port 8081):
POST /events body:{"type":"click-mouse","userId":"<string>","payload":{}} → 201 {"id":"<string>","type":"<string>","userId":"<string>","payload":{}};
POST /events body:{"type":"submit-form","userId":"<string>","payload":{}} → 201 (same shape);
GET /events → 200 JSON array containing previously posted events.

IAM listener (port 8082):
POST /users body:{"username":"<string>","password":"<string>"} → 201 {"id":"<string>","username":"<string>"} (no password field);
GET /users/{id} → 200 {"id":"<string>","username":"<string>"} or 404;
DELETE /users/{id} → 200/204 on existing user or 404 on missing;
POST /auth/login body:{"username":"<string>","password":"<string>"} → 200 {"token":"<string>"};
POST /auth/logout header:Authorization:Bearer <token> → 200/204;
POST /roles body:{"name":"<string>","permissions":["<string>"]} → 201 {"id":"<string>","name":"<string>","permissions":["<string>"]};
GET /roles → 200 JSON array;
POST /users/{id}/roles body:{"roleId":"<string>"} → 200/201;
GET /users/{id}/roles → 200 JSON array;
POST /authz/check body:{"userId":"<string>","permission":"<string>"} → 200 {"allowed":<bool>}.

## Context

### Level 1 — fbc6d9-0000-build-1


## Design

**Language:** Go. Module `github.com/cornjacket/platform-monolith`.

**Files to produce:**

### `cmd/platform/main.go`
```
package main

import (
    "context"
    "log"
    "net/http"
    "os/signal"
    "syscall"

    "github.com/cornjacket/platform-monolith/internal/iam"
    "github.com/cornjacket/platform-monolith/internal/metrics"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    store := metrics.NewEventStore()
    metricsSrv := &http.Server{Addr: ":8081", Handler: metrics.NewRouter(store)}
    iamSrv := &http.Server{Addr: ":8082", Handler: iam.NewMux()}

    go func() { log.Println(metricsSrv.ListenAndServe()) }()
    go func() { log.Println(iamSrv.ListenAndServe()) }()

    <-ctx.Done()
    stop()
    metricsSrv.Shutdown(context.Background())
    iamSrv.Shutdown(context.Background())
}
```

### `cmd/platform/platform_test.go`
```
package main_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    iampkg "github.com/cornjacket/platform-monolith/internal/iam"
    "github.com/cornjacket/platform-monolith/internal/metrics"
)

var (
    metricsBase string
    iamBase     string
)

func TestMain(m *testing.M) {
    mSrv := httptest.NewServer(metrics.NewRouter(metrics.NewEventStore()))
    iSrv := httptest.NewServer(iampkg.NewMux())
    metricsBase = mSrv.URL
    iamBase = iSrv.URL
    code := m.Run()
    mSrv.Close()
    iSrv.Close()
    os.Exit(code)
}
```

Individual test functions (`TestPostEventClickMouse`, `TestPostEventSubmitForm`, `TestGetEvents`, `TestPostUser`, `TestGetUser`, `TestDeleteUser`, `TestLogin`, `TestLogout`, `TestPostRole`, `TestGetRoles`, `TestAssignRole`, `TestGetUserRoles`, `TestAuthzCheck`) each use `net/http` against `metricsBase` or `iamBase`.

**Constraints:**
- `cmd/platform` is the sole `package main` in the repo — no other directory may declare `package main`.
- Tests use `package main_test` (external test package) so they do not add a second main package.
- `httptest.NewServer` provides real TCP listeners satisfying the "net/http against live listeners" requirement.
- `platform_test.go` imports `internal/iam` and `internal/metrics` directly (not via `cmd/platform`), so it builds cleanly without importing the main package itself.
- Tests that depend on prior state (e.g. `TestGetUser` needs a created user ID) chain within a single test function or use `TestMain`-level setup; they do not rely on test execution order across top-level `Test*` functions unless explicitly sequenced within one function.
- The IAM state (users, roles, tokens) is fresh per test run because `TestMain` creates a new `iam.NewMux()` — each `httptest.NewServer` wraps a brand-new in-memory store.

**Dependencies:** stdlib only (`net/http`, `net/http/httptest`, `encoding/json`, `os/signal`, `syscall`). No new external packages.

## Acceptance Criteria

1. `go build ./cmd/platform/...` succeeds with no errors.
2. `go vet ./...` reports no issues.
3. There is exactly one directory in the repo containing `package main` (`cmd/platform`); `grep -r "^package main" --include="*.go" .` returns only files under `cmd/platform/`.
4. POST `metricsBase/events` with `{"type":"click-mouse","userId":"u1","payload":{}}` returns HTTP 201 and a JSON body containing fields `id` (non-empty string), `type` ("click-mouse"), `userId` ("u1"), `payload` (object).
5. POST `metricsBase/events` with `{"type":"submit-form","userId":"u2","payload":{}}` returns HTTP 201 and a JSON body with `type` equal to "submit-form".
6. GET `metricsBase/events` returns HTTP 200 and a JSON array containing at least the two events created in criteria 4–5 (matched by their `id` fields).
7. POST `iamBase/users` with `{"username":"alice","password":"secret"}` returns HTTP 201 and a JSON body with fields `id` (non-empty string) and `username` ("alice"); the response body must NOT contain a `password` field.
8. GET `iamBase/users/{id}` using the ID from criterion 7 returns HTTP 200 with the same `id` and `username`; GET with a non-existent ID returns HTTP 404.
9. DELETE `iamBase/users/{id}` for an existing user returns HTTP 200 or 204; DELETE for a non-existent ID returns HTTP 404.
10. POST `iamBase/auth/login` with `{"username":"alice","password":"secret"}` (using a user created in criterion 7) returns HTTP 200 and a JSON body containing a non-empty `token` field.
11. POST `iamBase/auth/logout` with header `Authorization: Bearer <token>` (token from criterion 10) returns HTTP 200 or 204.
12. POST `iamBase/roles` with `{"name":"admin","permissions":["read","write"]}` returns HTTP 201 and a JSON body with fields `id` (non-empty string), `name` ("admin"), `permissions` (array containing "read" and "write").
13. GET `iamBase/roles` returns HTTP 200 and a JSON array containing at least the role created in criterion 12.
14. POST `iamBase/users/{id}/roles` with `{"roleId":"<roleId>"}` (using valid user and role IDs) returns HTTP 200 or 201.
15. GET `iamBase/users/{id}/roles` returns HTTP 200 and a JSON array containing the role assigned in criterion 14.
16. POST `iamBase/authz/check` with `{"userId":"<id>","permission":"read"}` returns HTTP 200 and `{"allowed":true}` for a user with the admin role (which has "read" permission); returns `{"allowed":false}` for a permission not held.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/platform-monolith/output && go test ./cmd/platform/...
```

