# main.go

Purpose: Entry point for the user service — constructs the store and handler, registers routes on a ServeMux, and starts the HTTP server on :8080.
All wiring of internal sub-packages happens here and nowhere else.

Tags: architecture, overview

## Responsibility

`main.go` is the composition root of the service. It owns the only site in the codebase where a concrete `*store.Store` is instantiated and passed to `handlers.New`, satisfying the `Store` interface declared by the handlers package. This is the single point of coupling between the two internal sub-packages.

## Functions

### `newMux() *http.ServeMux`

Constructs the full request-handling tree. It calls `store.New()`, passes the result to `handlers.New()`, and registers all routes via `Handler.RegisterRoutes`. The returned `*http.ServeMux` is ready to serve requests. This function is factored out of `main` to make it directly testable without starting a live server.

### `main()`

Calls `newMux()` and passes the result to `http.ListenAndServe(":8080", ...)`. Any startup failure (e.g. port in use) is fatal-logged and terminates the process. There is no graceful-shutdown path.

## Wiring Summary

| Constructed | Passed to | Purpose |
|-------------|-----------|---------|
| `*store.Store` (via `store.New()`) | `handlers.New()` | Provides the concrete Store implementation |
| `handlers.Handler` (via `handlers.New()`) | `Handler.RegisterRoutes(mux)` | Registers all `/users` routes |
| `*http.ServeMux` | `http.ListenAndServe` | Routes inbound requests to handlers |
