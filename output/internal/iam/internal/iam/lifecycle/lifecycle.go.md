Purpose: Documents the design and implementation decisions for the lifecycle package, which owns user CRUD and session-token management for the IAM listener.
Exposes a single public entry point, Handler(), that wires five HTTP routes onto a *http.ServeMux.

Tags: implementation, lifecycle

## Key types

### User
Holds an individual user record: `ID` (UUID string), `Username`, and `PasswordHash` (bcrypt, cost 12).
The password hash is never surfaced in HTTP responses.

### UserStore
Dual-indexed concurrent map: `byID` (id → *User) and `byUsername` (username → *User).
Both indexes are kept in sync on every insert and delete — there is no lazy sync step.
Protected by `sync.RWMutex`; reads use `RLock`, writes use `Lock`.

### TokenStore
Maps opaque token strings to userIDs.
Protected by its own `sync.RWMutex` independent of UserStore, preventing unnecessary contention.

## Main functions

### Handler() *http.ServeMux
Constructs zero-value stores, wires four `HandleFunc` registrations, and returns the mux.
The only exported symbol in the package.

### usersRoot (POST /users)
Validates non-empty username and password, checks for duplicate usernames under write-lock (read-then-check-then-insert is a single critical section to prevent TOCTOU races), hashes the password, inserts into both indexes, responds 201.

### usersID (GET or DELETE /users/{id})
ID extracted via `strings.TrimPrefix(r.URL.Path, "/users/")`.
DELETE acquires the UserStore write-lock, removes from both maps, then acquires the TokenStore write-lock and purges any tokens whose value equals the deleted userID (cascade invalidation). The two locks are released separately — never held simultaneously — to avoid deadlock.

### login (POST /auth/login)
Looks up user by username under RLock, compares bcrypt hash, generates a UUID token, stores it under TokenStore write-lock, responds 200 with token.

### logout (POST /auth/logout)
Parses `Authorization: Bearer <token>` header (returns 401 on missing prefix), acquires TokenStore write-lock, checks existence (returns 401 if absent), deletes and responds 200.

## Design decisions

- **Separate mutexes for UserStore and TokenStore**: avoids holding a coarse lock across both stores during DELETE cascade; each store is independently locked.
- **Write-lock for duplicate-username check in POST /users**: the check and insert must be atomic to prevent two concurrent requests from both passing the check and inserting the same username.
- **bcrypt cost 12**: balances security against test latency; cost is hardcoded per the design spec.
- **No exported store types**: stores are package-private; only `Handler()` is exported, keeping the public surface minimal.
- **405 for wrong methods**: each handler returns `405 Method Not Allowed` rather than `404` for unsupported HTTP verbs, which is more correct per HTTP semantics.
