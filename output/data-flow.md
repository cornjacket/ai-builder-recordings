# data-flow

Purpose: End-to-end description of how the user service starts up and processes HTTP requests — from `main.go` wiring through the internal package boundary to the JSON response returned to the client.

Tags: data-flow, architecture

## Startup Sequence

Before any request can be served, `main.go` performs one-time wiring. The diagram below shows the construction order and how the resulting objects relate to each other:

```
main.go
  |
  |-- store.New() ──────────────────────► *store.Store
  |                                          |
  |-- handlers.New(*store.Store) ──────► handlers.Handler
  |                                          |
  |-- Handler.RegisterRoutes(*ServeMux) ─► *http.ServeMux
  |
  └── http.ListenAndServe(":8080", mux)   [blocks]
```

The `*store.Store` is created first because `handlers.New` requires a value satisfying the `Store` interface. Once the mux is built, the process blocks in `ListenAndServe` until a fatal error occurs.

## Request Path

Once the server is running, every inbound request follows the same path:

```
HTTP Client
    |
    | HTTP request (e.g. DELETE /users/{id})
    v
*http.ServeMux  (main.go)
    |
    | matched route → Handler method
    v
handlers.Handler  [internal/userservice/handlers]
    |
    | calls Store interface method (Create / Get / Update / Delete)
    v
store.Store  [internal/userservice/store]
    |
    | mutex-guarded map[string]User operation
    | returns (User, error) or error
    v
handlers.Handler
    |
    | encodes result as JSON (writeJSON / writeError)
    v
HTTP Client  (JSON response + status code)
```

## Interface Decoupling

`main.go` is the only file that imports both `handlers` and `store` directly. The handlers package declares the `Store` interface it needs; the store package provides the concrete type that satisfies it. Neither internal sub-package imports the other. This means both can be tested independently, and `newMux()` is the sole wiring site for the entire dependency graph.

## Concurrency Model

The store uses an `sync.RWMutex`: read operations hold a shared lock, write operations hold an exclusive lock. The handlers layer is stateless between requests. The `ServeMux` dispatches each request in its own goroutine (standard `net/http` behaviour), so concurrent requests are safe without any additional synchronisation in `main.go`.
