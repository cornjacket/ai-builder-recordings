# lifecycle.go

Purpose: Implements user CRUD and session token management for the IAM lifecycle sub-package,
exposing four HTTP endpoints via a single Handler backed by an in-memory store.
Passwords are hashed with bcrypt; session tokens are opaque UUIDs.

Tags: api, data-access, model

## Public API

### Types

**`User`**
Domain model for a registered user.

| Field | Type | Description |
|---|---|---|
| `ID` | `string` | UUID assigned at creation |
| `Username` | `string` | Unique display name |
| `PasswordHash` | `string` | bcrypt hash of the user's password |

**`Token`**
Represents an active session.

| Field | Type | Description |
|---|---|---|
| `Token` | `string` | UUID used as the bearer token |
| `UserID` | `string` | ID of the owning user |

**`Store`**
Holds all in-memory lifecycle state behind a single `sync.RWMutex`. Unexported; obtained only via `New()`.

**`Handler`**
Wraps a `*Store` and satisfies the lifecycle HTTP surface. Unexported fields only; obtained via `New()`.

### Functions and Methods

**`New() *Handler`**
Constructs a Handler with a freshly initialised Store. No configuration is required.

**`(h *Handler) RegisterRoutes(mux *http.ServeMux)`**
Registers the four lifecycle endpoints onto `mux`.

| Method | Path | Description |
|---|---|---|
| `POST` | `/users` | Create a new user; returns `{"id","username"}` |
| `GET` | `/users/{id}` | Fetch a user by ID; returns `{"id","username"}` |
| `DELETE` | `/users/{id}` | Delete a user by ID |
| `POST` | `/auth/login` | Validate credentials; returns `{"token"}` |
| `POST` | `/auth/logout` | Invalidate the bearer token in the `Authorization` header |

## Key Internals

**`Store` concurrency model**
All mutating operations (`addUser`, `deleteUser`, `addToken`, `deleteToken`) acquire a full write lock. Read operations (`getUser`, `getUserByName`, `getToken`) acquire a read lock. The dual-index design (`users` map keyed by ID, `byName` keyed by username) allows O(1) lookups by either key at the cost of keeping both maps consistent on delete.

**`addUser` duplicate detection**
Checks `byName` under a write lock before inserting. Returns `errUsernameTaken` (a typed `errDuplicate` string) on collision; the handler maps this to HTTP 409 Conflict.

**`handleLogin` credential check**
Resolves the username to a `User` via `getUserByName`, then calls `bcrypt.CompareHashAndPassword`. Both a missing user and a wrong password return HTTP 401 to avoid username enumeration.

**`handleLogout` token extraction**
Strips the `Bearer ` prefix from the `Authorization` header. If the prefix is absent or the resulting token string is not found in the store, the request is rejected with HTTP 401.

**`errDuplicate`**
A named `string` type that implements `error`. Used exclusively for the username-already-taken sentinel so callers can distinguish duplicate errors from other failures.

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/google/uuid` | Generates UUID strings for user IDs and session tokens |
| `golang.org/x/crypto/bcrypt` | Hashes passwords at `bcrypt.DefaultCost` and verifies them on login |
