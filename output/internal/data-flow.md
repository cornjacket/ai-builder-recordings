# data-flow

Purpose: Describes how the iam and metrics services are constructed and how HTTP requests flow through each at runtime.
The two services share no Go-level coupling; coordination happens exclusively in cmd/platform.

Tags: data-flow, architecture

## Construction (startup)

`cmd/platform` calls each service's `New()` constructor and mounts the returned `http.Handler` on its own `http.Server`:

```
cmd/platform main()
  │
  ├─► internal/iam.New()
  │     ├─► lifecycle.New()   → *lifecycle.Handler (owns user store + token store)
  │     ├─► authz.New()       → *authz.Handler    (owns role store + assignment store)
  │     └─► register both onto *http.ServeMux
  │                                │
  │                                └─► http.Handler bound to :8082
  │
  └─► internal/metrics.New()
        ├─► store.New()        → *store.Store (owns in-memory event slice)
        ├─► handlers.New(store) → *handlers.Handler
        └─► register onto *http.ServeMux
                                 │
                                 └─► http.Handler bound to :8081
```

Each `New()` returns an opaque `http.Handler`; neither service holds a reference to the other.

## Request-time paths

### iam service (:8082)

```
HTTP client
    │
    ▼
iam ServeMux
    │
    ├─── /lifecycle/* ──► lifecycle.Handler
    │                         │
    │         ┌───────────────┴──────────────────┐
    │         ▼                                   ▼
    │   POST /lifecycle/users            POST /lifecycle/login
    │   (addUser → store.users)          (handleLogin → store.tokens)
    │                                             │
    │                                    POST /lifecycle/logout
    │                                    (handleLogout → store.tokens)
    │
    └─── /authz/* ──────► authz.Handler
                              │
              ┌───────────────┼────────────────────┐
              ▼               ▼                    ▼
        POST /authz/roles  POST /authz/assign  GET /authz/check
        (addRole)          (assignRole)        (hasPermission)
```

A user ID issued by `lifecycle` is the shared key that a client passes to `authz` endpoints — there is no direct Go call between the two sub-packages; the client is the integration point.

### metrics service (:8081)

```
HTTP client
    │
    ▼
metrics ServeMux
    │
    ├─── POST /events ──► handlers.PostEvents
    │                         │
    │                   validate type allowlist
    │                         │
    │                         ▼
    │                   store.Add(event)   → appends to in-memory slice
    │
    └─── GET  /events ──► handlers.GetEvents
                              │
                              ▼
                        store.List()       → returns slice copy as JSON
```

## Cross-service isolation

`iam` and `metrics` have no shared state, no shared interfaces, and no import relationship. They are isolated services that happen to live under the same `internal/` directory. The only entity that knows about both is `cmd/platform`, which constructs and serves them independently.
