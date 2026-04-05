# data-flow

Purpose: Describes how an HTTP request moves through the userservice composite — from the ServeMux registered in main.go, through the handlers layer, into the store, and back as a JSON response.

Tags: data-flow, architecture

## Request Path

Every request to the user service follows the same linear path: the `ServeMux` in `main.go` dispatches to a `Handler` method in `handlers`, which calls one method on the `Store` interface, and the result is encoded as JSON and written to the response.

```
HTTP Client
    |
    | HTTP request (e.g. POST /users)
    v
main.go — *http.ServeMux
    |
    | matched route dispatched to Handler method
    v
handlers.Handler  (reads/writes http.ResponseWriter, *http.Request)
    |
    | calls Store interface method (Create / Get / Update / Delete)
    v
store.Store  (mutex-guarded map[string]User)
    |
    | returns (User, error) or error
    v
handlers.Handler  (encodes result as JSON via writeJSON / writeError)
    |
    | HTTP response (JSON body, status code)
    v
HTTP Client
```

## Interface Boundary

The coupling between `handlers` and `store` is intentionally indirect. `handlers` defines a `Store` interface; `store` provides the concrete `*store.Store` type that satisfies it. `main.go` is the only site where the concrete type is referenced — it passes a `*store.Store` to `handlers.New`. This means the two sub-packages are independently testable with no import cycle.

## Concurrency

All store operations are safe for concurrent use. Read operations (`Get`) hold a shared read lock; write operations (`Create`, `Update`, `Delete`) hold an exclusive write lock. The `handlers` package is stateless between requests and adds no additional synchronisation of its own.
