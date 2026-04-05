# Theory of Operation — iam

Purpose: Describes the data-flow and component interactions within the IAM composite, showing how lifecycle and authz sub-packages collaborate to serve all ten IAM endpoints.

Tags: architecture, design

## Overview

The `integrate-iam` wiring step registers handlers from both `lifecycle` and `authz` on a single `http.ServeMux`, then passes that mux to the platform's HTTP server for port 8082. Neither sub-package knows about the other at compile time; they share no types and communicate only through the HTTP layer.

## Data-Flow

```
HTTP Request (port 8082)
        │
        ▼
  http.ServeMux (iam mux)
        │
        ├─── /users, /auth/*  ──────────────────► lifecycle package
        │                                              │
        │   ┌──────────────────────────────────────┐  │
        │   │  UserStore (map[id]User)              │◄─┘
        │   │  TokenStore (map[token]userID)        │
        │   └──────────────────────────────────────┘
        │
        └─── /roles, /users/{id}/roles, /authz/*  ─► authz package
                                                        │
             ┌──────────────────────────────────────┐  │
             │  RoleStore  (map[id]Role)             │◄─┘
             │  UserRoles  (map[userID][]roleID)     │
             └──────────────────────────────────────┘
```

## Request Lifecycle

```
Incoming request
      │
      ▼
 ServeMux.ServeHTTP
      │
      ├─ prefix /users or /auth  →  lifecycle.Handler.ServeHTTP
      │         │
      │         ├─ POST /users        →  register; write UserStore
      │         ├─ GET  /users/{id}   →  read UserStore
      │         ├─ DELETE /users/{id} →  delete from UserStore
      │         ├─ POST /auth/login   →  verify UserStore; issue token; write TokenStore
      │         └─ POST /auth/logout  →  delete from TokenStore
      │
      └─ prefix /roles or /authz  →  authz.Handler.ServeHTTP
                │
                ├─ POST /roles              →  write RoleStore
                ├─ GET  /roles              →  read RoleStore
                ├─ POST /users/{id}/roles   →  write UserRoles
                ├─ GET  /users/{id}/roles   →  read UserRoles + join RoleStore
                └─ POST /authz/check        →  read UserRoles + RoleStore; eval permission
```

## State

Each sub-package owns its state exclusively:

| Sub-package | State | Shared? |
|-------------|-------|---------|
| lifecycle | UserStore, TokenStore | No |
| authz | RoleStore, UserRoles | No |

Cross-cutting concerns (e.g. validating a user ID exists before assigning a role) are not required by the acceptance spec and are deliberately omitted to keep the packages independent.
