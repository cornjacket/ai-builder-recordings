# iam.go

Purpose: Entry-point wiring file for the IAM composite — `New()` instantiates both sub-package handlers and registers all ten HTTP routes on a single `ServeMux`.

Tags: architecture

## Public API

`New() http.Handler` is the only exported symbol. It constructs a fresh `http.ServeMux`, delegates route registration to both sub-packages in sequence, and returns the mux as an `http.Handler`. The caller (`cmd/platform/main.go`) binds the returned handler to `:8082`.

## Wiring Sequence

Neither sub-package knows about the other; `iam.go` is the only place they are composed.

1. `lifecycle.New()` — allocates the lifecycle `Handler` (user store + token store, both empty).
2. `lifecycle.Handler.RegisterRoutes(mux)` — registers five routes under `/users` and `/auth`.
3. `authz.New()` — allocates the authz `Handler` (role store + assignment store, both empty).
4. `authz.Handler.RegisterRoutes(mux)` — registers five routes under `/roles`, `/users/{id}/roles`, and `/authz/check`.
5. Returns the mux as `http.Handler`.

## Dependencies

| Import | Role |
|--------|------|
| `net/http` | `ServeMux` and `http.Handler` interface |
| `internal/iam/lifecycle` | User lifecycle handler |
| `internal/iam/authz` | RBAC handler |
