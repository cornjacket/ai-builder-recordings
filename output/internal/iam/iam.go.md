Purpose: Wires lifecycle and authz sub-muxes into the single IAM *http.ServeMux returned by NewMux().
Callers pass this mux directly to http.ListenAndServe — no port binding occurs here.

Tags: implementation, iam

## Key functions

### NewMux() *http.ServeMux

Calls `lifecycle.Handler()` and `authz.Handler()` to obtain two independent sub-muxes, then registers
routes on a fresh top-level mux:

| Pattern | Handler |
|---|---|
| `/users` | lifecycle |
| `/auth/login` | lifecycle |
| `/auth/logout` | lifecycle |
| `/roles` | authz |
| `/authz/check` | authz |
| `/users/` (catch-all) | custom HandlerFunc (see below) |

## Design decisions

### /users/ catch-all with suffix check

Go's `net/http` mux treats `/users/` as a subtree pattern that matches any path beginning with
`/users/`.  A request for `/users/{id}` must reach lifecycle, while `/users/{id}/roles` must reach
authz.  The disambiguator is a single `strings.HasSuffix(r.URL.Path, "/roles")` check.

This is safe because the only valid sub-resource under `/users/{id}` is `/roles`; no other valid path
ends in `/roles`.

### Unmodified requests

Both sub-muxes receive the original `*http.Request` with its full `URL.Path` intact.  Each sub-mux's
internal path extraction (`strings.TrimPrefix`, `strings.Split`) was written to operate on the full
path, so no stripping is needed at this layer.

### No global state

`NewMux()` calls `Handler()` on each sub-package on every invocation, ensuring each returned mux has
its own independent in-memory stores.  Tests can call `NewMux()` multiple times without cross-test
interference.
