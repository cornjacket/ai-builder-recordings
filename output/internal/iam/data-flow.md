# data-flow.md

Purpose: Describes how data moves between the three components of the IAM composite — lifecycle, authz, and the iam.go wiring layer.

Tags: data-flow

## Overview

The IAM composite is assembled from two independent sub-packages wired together by a thin top-level file. Neither sub-package imports the other; all coupling is via shared string user IDs on the HTTP surface.

## Startup Wiring

```
main.go (:8082)
    └── iam.New()
            ├── lifecycle.New()  →  lifecycle.Handler
            │       └── RegisterRoutes(mux)  →  5 routes mounted
            └── authz.New()      →  authz.Handler
                    └── RegisterRoutes(mux)  →  5 routes mounted
                                    ↓
                            http.Handler returned to main.go
```

`New()` calls each sub-package constructor exactly once, registers their routes on the same `ServeMux`, then returns the mux. The two in-memory stores are initialised empty and never share a reference.

## Request-Time Data Flow

```
HTTP client
    │
    ▼
http.ServeMux (iam.go)
    │
    ├── /users, /auth/*  ──────────────────►  lifecycle.Handler
    │                                              │
    │                                    User store (map[id]User)
    │                                    Token store (map[token]userID)
    │                                              │
    │                                        JSON response
    │
    └── /roles, /users/{id}/roles,
        /authz/check  ─────────────────────►  authz.Handler
                                                   │
                                         Role store (map[id]Role)
                                         Assignment store (map[userID][]roleID)
                                                   │
                                             JSON response
```

### Cross-component coupling

`authz` references users only by the opaque string `userID` carried in request paths (e.g. `POST /users/{id}/roles`). It never queries `lifecycle` directly. The client is responsible for obtaining a user ID from `lifecycle` (via `POST /users`) and supplying it to `authz` endpoints. This is the only cross-component data path, and it runs entirely through the HTTP surface.

```
Client workflow example
───────────────────────
  1. POST /users          → lifecycle → returns { id, username }
  2. POST /roles          → authz    → returns { id, name, permissions }
  3. POST /users/{id}/roles → authz  → assigns role to user (string ID only)
  4. POST /authz/check    → authz    → resolves permission via assigned roles
```

## State Isolation

Each sub-package owns its state exclusively; there are no shared data structures, mutexes, or channels between them. A failure or inconsistency in one store does not affect the other.
