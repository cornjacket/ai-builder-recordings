Purpose: Documents the design and key decisions of the authz package's role store and RBAC handler.
Covers store structure, mutex strategy, routing approach, and endpoint behaviour.

Tags: implementation, authz

## Key Types

**`Role`** — data model for an access-control role; serialised to/from JSON on all endpoints.

**`handler`** — owns both in-memory stores (`roleStore` and `userRoles`) plus the single `sync.RWMutex` that guards them. One struct, one lock — keeps the ownership model simple.

## Main Functions

**`Handler()`** — sole public entry point. Allocates the handler, registers three path patterns on a new `*http.ServeMux`, and returns it. Callers (integrate-iam) mount the mux under their own prefix if needed.

**`handleRoles`** — dispatches `POST /roles` (create) and `GET /roles` (list) based on `r.Method`.

**`handleUsers`** — handles `/users/{id}/roles` for both POST and GET. Splits `r.URL.Path` on `/` to extract the user ID; returns 404 for malformed paths.

**`handleAuthzCheck`** — walks the user's assigned role IDs, fetches each from the store, and scans permissions for an exact string match. Always returns 200 with `{"allowed":bool}`; unknown users yield `false` without error.

## Design Decisions

**Single mutex for both maps** — both stores are modified together only during `assignRole` (which reads roleStore then writes userRoles). A single `sync.RWMutex` is simpler and avoids lock-ordering bugs that two separate mutexes would introduce.

**Write lock held across read+modify in `assignRole`** — the check for roleID existence and the append to `userRoles` must be atomic. If they were split, a concurrent role deletion could cause a stale reference to be stored. Holding the full `Lock` throughout keeps the operation consistent.

**`/users/` trailing-slash prefix match** — Go 1.21's `http.ServeMux` does not support method or wildcard patterns. A trailing-slash match is the standard way to capture an arbitrary sub-path with one `HandleFunc` call; the handler then parses the path itself.

**Idempotent role assignment** — before appending a roleID, `assignRole` scans the existing slice. For the expected cardinality (single-digit roles per user), a linear scan is adequate and avoids introducing a set dependency.

**`POST /authz/check` never returns 404** — an unknown `userId` simply has no roles, so the permission walk yields `false`. Returning 404 would leak whether a user exists; 200 `{"allowed":false}` is consistent with a deny-by-default policy.
