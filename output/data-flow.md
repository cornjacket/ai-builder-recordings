# data-flow.md

Purpose: Cross-component data-flow for the platform monolith — how cmd/platform constructs both services at startup and how HTTP requests are dispatched to each at runtime.
Covers the full path from process start through service initialisation to request handling across both ports.

Tags: data-flow, architecture

## Startup: Object-Graph Construction

At process start, `cmd/platform/main.go` calls each service's `New()` constructor and then starts two goroutines, one `http.ListenAndServe` per handler.

```
main()
  │
  ├─── internal/iam.New()
  │         │
  │         ├─── lifecycle.New()   → *lifecycle.Store (users, tokens in memory)
  │         │         └─── *lifecycle.Handler (wraps store)
  │         ├─── authz.New()       → *authz.Store (roles, assignments in memory)
  │         │         └─── *authz.Handler (wraps store)
  │         └─── ServeMux          ← routes registered for both handlers
  │                   └─── returns http.Handler  ──► :8082
  │
  └─── internal/metrics.New()
            │
            ├─── store.New()       → *store.Store (events slice in memory)
            ├─── handlers.New(store) → *handlers.Handler (wraps store via Storer interface)
            └─── ServeMux          ← POST /events, GET /events registered
                      └─── returns http.Handler  ──► :8081

goroutine A: http.ListenAndServe(":8082", iamHandler)
goroutine B: http.ListenAndServe(":8081", metricsHandler)
main goroutine: blocks on select{}
```

Neither `internal/iam` nor `internal/metrics` imports the other. All composition happens in `main()`.

## Request-Time Routing

### IAM service (:8082)

```
HTTP client
    │
    ▼
:8082  iam ServeMux
    │
    ├── /users, /users/{id}          → lifecycle.Handler
    │       reads/writes lifecycle.Store (RWMutex-guarded users map, tokens map)
    │
    ├── /auth/login, /auth/logout    → lifecycle.Handler
    │       reads/writes lifecycle.Store
    │
    ├── /roles, /users/{id}/roles    → authz.Handler
    │       reads/writes authz.Store (RWMutex-guarded roles map, assignments map)
    │
    └── /authz/check                 → authz.Handler
            reads authz.Store (permission evaluation — no writes)
```

### Metrics service (:8081)

```
HTTP client
    │
    ▼
:8081  metrics ServeMux
    │
    ├── POST /events  → handlers.Handler
    │       validates type allowlist (click-mouse | submit-form)
    │       calls store.Storer.Add() → appends to store.Store slice (Mutex-guarded)
    │       responds 201 with stored Event (UUID assigned by store)
    │
    └── GET /events   → handlers.Handler
            calls store.Storer.List() → returns copy of events slice
            responds 200 with JSON array (empty array [] when no events)
```

## Cross-Service Isolation

The two services share no Go-level state. The only cross-service coupling is at the HTTP client level: a caller that creates a user via `POST /users` (iam) and then assigns a role via `POST /users/{id}/roles` (iam/authz) uses the same user ID string, but this coordination is entirely the client's responsibility. The metrics service is independent of both lifecycle and authz.

Both services use in-process memory only. State is lost when the process exits.
