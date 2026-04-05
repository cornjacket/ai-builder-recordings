# metrics.go.md

Purpose: Documents the package-level New() constructor that wires the store and handlers sub-packages into a single http.Handler. This file is the only source at the metrics/ level; all domain logic lives in the sub-packages.

Tags: architecture, data-flow

## Public API

```go
func New() http.Handler
```

Returns a configured `*http.ServeMux` with a single route registered. The caller is responsible for wrapping the returned handler in an `http.Server` (done by `cmd/platform/main.go` on `:8081`).

## Wiring Sequence

`New()` performs three steps in order:

1. **Store construction** — calls `store.New()`, which allocates an in-memory `Store` with an initialised `sync.RWMutex` and an empty event slice.
2. **Handler construction** — calls `handlers.New(s)`, injecting the store as a `store.Storer` interface. No import cycle is created because `handlers` depends on the `store` package, not on this parent package.
3. **Route registration** — registers a single catch-all handler on `/events` that dispatches to `h.PostEvents` (POST), `h.GetEvents` (GET), or returns `405 Method Not Allowed` for all other methods.

## Route Table

| Method | Path | Handler |
|--------|------|---------|
| POST | `/events` | `handlers.Handler.PostEvents` |
| GET | `/events` | `handlers.Handler.GetEvents` |
| * | `/events` | `405 Method Not Allowed` |

## Dependencies

| Import | Role |
|--------|------|
| `net/http` | `http.Handler`, `http.ServeMux`, `http.MethodPost`, `http.MethodGet`, `http.Error` |
| `internal/metrics/handlers` | HTTP layer; provides `New(store.Storer)` and the `PostEvents`/`GetEvents` methods |
| `internal/metrics/store` | Data layer; provides `New()` and the `Storer` interface |

## Design Notes

The method dispatch switch in `New()` is the only logic that lives at the `metrics` package level. By keeping route multiplexing here rather than inside `handlers`, the handlers package stays free of HTTP method awareness and can be tested with a direct function call.
