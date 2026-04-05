# lifecycle

## Goal

In-memory user store and token-based authentication handler for the IAM listener.

Endpoints:
POST /users {"username":string,"password":string} → 201 {"id":string,"username":string} or 400;
GET /users/{id} → 200 {"id":string,"username":string} or 404;
DELETE /users/{id} → 200 or 404;
POST /auth/login {"username":string,"password":string} → 200 {"token":string} or 401;
POST /auth/logout header Authorization:Bearer <token> → 200 or 401.

In-memory stores: UserStore map[id]→{id string, username string, passwordHash string}, TokenStore map[token]→userID string. IDs generated with crypto/rand or math/rand UUID. Passwords hashed with bcrypt or sha256. Exposes an http.Handler (or *http.ServeMux) for registration by integrate-iam.

## Context

### Level 1 — fbc6d9-0000-build-1

### Level 2 — fbc6d9-0001-iam
Identity and access management listener (port 8082) composed of two sub-components — lifecycle (user CRUD and authentication) and authz-rbac (roles and permission checks) — that together handle all ten IAM HTTP endpoints.

## Design

**Language:** Go 1.21  
**Package:** `lifecycle`  
**Output directory:** `output/internal/iam/internal/iam/lifecycle/`

### Files

| File | Responsibility |
|------|----------------|
| `lifecycle.go` | Stores, handler struct, `Handler()`, all five HTTP handler methods |
| `lifecycle_test.go` | Table-driven tests for all five endpoints |

### Data models

```go
type User struct {
    ID           string
    Username     string
    PasswordHash string
}

type UserStore struct {
    mu         sync.RWMutex
    byID       map[string]*User   // id → *User
    byUsername map[string]*User   // username → *User
}

type TokenStore struct {
    mu     sync.RWMutex
    tokens map[string]string // token → userID
}
```

### Public API

```go
// Handler constructs both stores and returns a *http.ServeMux with all five routes registered.
func Handler() *http.ServeMux
```

### Route registration (Go 1.21 stdlib)

```go
mux := http.NewServeMux()
mux.HandleFunc("/users",   h.usersRoot)  // POST
mux.HandleFunc("/users/",  h.usersID)    // GET, DELETE — trailing slash catches /users/{id}
mux.HandleFunc("/auth/login",  h.login)
mux.HandleFunc("/auth/logout", h.logout)
```

Each handler switches on `r.Method` and returns `405 Method Not Allowed` for unexpected methods. The ID in `/users/{id}` is extracted as `strings.TrimPrefix(r.URL.Path, "/users/")`.

### Handler logic

**POST /users (`usersRoot`):**
- Decode `{username, password}` from body; return 400 on decode error or empty fields.
- Lock UserStore write; check byUsername for duplicate; return 400 if found.
- Hash password with `bcrypt.GenerateFromPassword([]byte(password), 12)`.
- Generate ID with `uuid.NewString()`.
- Insert into both maps.
- Write 201 `{"id":..., "username":...}`.

**GET /users/{id} (`usersID` with method GET):**
- Extract ID from path.
- RLock UserStore; lookup byID; return 404 if missing.
- Write 200 `{"id":..., "username":...}`.

**DELETE /users/{id} (`usersID` with method DELETE):**
- Extract ID from path.
- Write-lock UserStore; lookup byID; return 404 if missing.
- Delete from both maps.
- Write-lock TokenStore; iterate tokens and delete any whose value matches the userID.
- Write 200 empty body.

**POST /auth/login (`login`):**
- Decode `{username, password}`.
- RLock UserStore; lookup byUsername; return 401 if not found.
- `bcrypt.CompareHashAndPassword`; return 401 on mismatch.
- Generate token with `uuid.NewString()`.
- Write-lock TokenStore; store `token → userID`.
- Write 200 `{"token":...}`.

**POST /auth/logout (`logout`):**
- Read `Authorization` header; parse `Bearer <token>`; return 401 if missing or malformed.
- Write-lock TokenStore; check presence; return 401 if not found.
- Delete token.
- Write 200 empty body.

### Dependencies

- `github.com/google/uuid` — already in go.mod
- `golang.org/x/crypto/bcrypt` — IMPLEMENTOR must run `go get golang.org/x/crypto` and update go.sum
- Standard library: `encoding/json`, `net/http`, `strings`, `sync`

### Non-obvious constraints

- `byUsername` map is kept in sync with `byID` on every insert and delete — no separate index sync step.
- DELETE cascades token invalidation: iterate TokenStore and remove any token whose value equals the deleted userID.
- All JSON error responses use the shape `{"error":"<message>"}`.
- `Handler()` is the only exported symbol; stores are not exported.

## Acceptance Criteria

1. `POST /users` with `{"username":"alice","password":"secret"}` returns HTTP 201 and a JSON body containing `"id"` (non-empty string) and `"username":"alice"`; the body must not contain any password field.
2. `POST /users` with a missing `username` field or missing `password` field returns HTTP 400.
3. `POST /users` called twice with the same username returns HTTP 400 on the second call.
4. `GET /users/{id}` using the ID returned from a successful `POST /users` returns HTTP 200 with `{"id":"<same-id>","username":"alice"}`.
5. `GET /users/{id}` with an ID that does not exist returns HTTP 404.
6. `DELETE /users/{id}` for an existing user returns HTTP 200.
7. `DELETE /users/{id}` for an ID that does not exist returns HTTP 404.
8. After `DELETE /users/{id}`, a subsequent `GET /users/{id}` for the same ID returns HTTP 404.
9. `POST /auth/login` with correct username and password returns HTTP 200 and a JSON body containing `"token"` (non-empty string).
10. `POST /auth/login` with the correct username but wrong password returns HTTP 401.
11. `POST /auth/login` with a username that was never registered returns HTTP 401.
12. `POST /auth/logout` with `Authorization: Bearer <token>` where the token was issued by a successful login returns HTTP 200.
13. `POST /auth/logout` called a second time with the same (now-invalidated) token returns HTTP 401.
14. `POST /auth/logout` with no `Authorization` header returns HTTP 401.
15. `POST /auth/logout` with `Authorization: Token <token>` (wrong scheme) returns HTTP 401.
16. After `DELETE /users/{id}`, a `POST /auth/logout` using a token that was issued for that user before deletion returns HTTP 401 (token was cascade-invalidated).
17. Running `go test -race ./...` from the module root targeting this package produces no data-race warnings when concurrent requests are made to the handler.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/platform-monolith/output && go test -race ./internal/iam/internal/iam/lifecycle/...
```

