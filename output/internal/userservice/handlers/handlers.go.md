Purpose: HTTP CRUD handlers for the user management API, wired to a Store interface.
Routes are registered on a net/http.ServeMux using Go 1.22+ method-qualified patterns.

Tags: implementation, handlers

## Key types

| Type | Responsibility |
|---|---|
| `Store` | Interface the handlers depend on; satisfied structurally by `*store.Store` |
| `createRequest` / `updateRequest` | Decode incoming JSON request bodies |
| `userResponse` | Encode outgoing JSON responses with `id`, `name`, `email` |

## Main functions

- **`New(s Store) http.Handler`** — creates a `ServeMux`, registers the four routes, returns it.
- **`handleCreate`** — decodes body → calls `Store.Create` → 201 + JSON.
- **`handleGet`** — extracts `{id}` via `r.PathValue` → `Store.Get` → 200 + JSON or 404.
- **`handleUpdate`** — extracts `{id}`, decodes body → `Store.Update` → 200 + JSON or 404.
- **`handleDelete`** — extracts `{id}` → `Store.Delete` → 204 (no body) or 404.
- **`writeNotFound`** — sets `Content-Type: application/json`, writes 404 + `{}`.

## Design decisions

- `Content-Type` header is set **before** `WriteHeader` on success paths so it is included in the response. `writeNotFound` follows the same ordering.
- `json.NewEncoder(w).Encode(v)` appends a trailing newline, which is acceptable per the design spec.
- The `Store` interface is defined in this package (not imported) so the handlers package does not depend on the concrete `*store.Store` type — only on `store.User` for the return values.
- Go 1.22+ `ServeMux` pattern syntax (`"POST /users"`) is required; the minimum version is enforced by `go.mod`.
