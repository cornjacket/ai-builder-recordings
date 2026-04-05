# authz.go

Purpose: Implements the RBAC sub-package — defines role and assignment models, an in-memory store protected by a read/write mutex, and an HTTP Handler that registers five permission-management routes.

Tags: api, model

## Public API

### Types

**`Role`**
```go
type Role struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Permissions []string `json:"permissions"`
}
```
Domain model for a named role carrying a list of permission strings. IDs are assigned by the store at creation time.

**`UserRoleAssignment`**
```go
type UserRoleAssignment struct {
    UserID string
    RoleID string
}
```
Records that a specific user has been granted a specific role. Assignments are append-only; there is no revocation path.

**`Handler`**
```go
type Handler struct{ s *store }
```
Public entry point for the authz sub-package. Holds a reference to the unexported in-memory store.

### Functions

**`New() *Handler`**
Allocates and returns a Handler backed by an initialised store. This is the only way to obtain a Handler.

**`(h *Handler) RegisterRoutes(mux *http.ServeMux)`**
Registers the five authz routes on the provided mux:

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/roles` | Create a new role |
| `GET` | `/roles` | List all roles |
| `POST` | `/users/{id}/roles` | Assign a role to a user |
| `GET` | `/users/{id}/roles` | List roles held by a user |
| `POST` | `/authz/check` | Check whether a user has a named permission |

## Key Internals

**`store`** — unexported struct holding `sync.RWMutex`, a `map[string]Role`, and a slice of `UserRoleAssignment`. Read operations acquire a read lock; writes acquire the exclusive lock.

**`store.assignRole(userID, roleID) bool`** — validates that the target role exists before appending the assignment; returns `false` (surfaced as HTTP 404) if the role is unknown.

**`store.hasPermission(userID, permission) bool`** — iterates all assignments for the user, resolves each to its Role, and scans the `Permissions` slice for a string match. O(n) over assignments × permissions; acceptable for the expected scale.

**`store.addRole`** — normalises a `nil` permissions argument to an empty slice before storing, so JSON serialisation never emits `null` for the permissions field.

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/google/uuid` | Generates unique IDs for newly created roles |
