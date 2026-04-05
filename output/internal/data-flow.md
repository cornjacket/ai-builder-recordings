# data-flow

Purpose: Describes how an inbound HTTP request crosses the internal/ package boundary and is processed by the userservice sub-tree before a response is returned to the caller.
Covers the entry point in main.go and the handoff into the internal sub-packages.

Tags: data-flow, architecture

## Entry into the Internal Boundary

All execution that uses the `internal/` packages originates in `main.go`. At startup, `main.go` constructs a `*store.Store` and a `handlers.Handler`, registers routes on a `*http.ServeMux`, and begins listening on `:8080`. Once the server is running, the only interaction with the internal packages is inbound HTTP requests dispatched by the mux.

```
main.go (project root)
    |
    | instantiates store.New() → *store.Store
    | instantiates handlers.New(*store.Store) → handlers.Handler
    | registers routes via Handler.RegisterRoutes(*http.ServeMux)
    |
    v
[ internal/ boundary ]
    |
    | inbound HTTP request (e.g. PUT /users/{id})
    v
userservice/handlers  — decodes request, calls Store interface method
    |
    v
userservice/store     — executes mutex-guarded map operation, returns result
    |
    v
userservice/handlers  — encodes result as JSON, writes HTTP response
    |
    v
[ response exits internal/ boundary → HTTP client ]
```

## Package Coupling

The two sub-packages inside `userservice/` are decoupled from each other at the import level. `handlers` defines the `Store` interface it requires; `store` provides the concrete type that satisfies it. Neither package imports the other. Only `main.go`, outside the `internal/` boundary, holds the concrete type reference and performs the wiring. This structure keeps each sub-package independently testable.
