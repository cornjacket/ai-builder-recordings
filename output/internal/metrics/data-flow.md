# data-flow.md

Purpose: Describes how data moves between the store and handlers sub-packages within the metrics composite, covering both startup wiring and request-time paths.

Tags: data-flow, architecture

## Startup Wiring

When `metrics.New()` is called by the binary, it constructs the object graph bottom-up before any HTTP traffic arrives.

```
metrics.New()
    │
    ├─► store.New()
    │       └─► *store.Store   (empty []Event, initialised RWMutex)
    │               │
    ├─► handlers.New(s)        (s passed as store.Storer interface)
    │       └─► *handlers.Handler
    │
    └─► http.NewServeMux()
            └─► mux.HandleFunc("/events", dispatch)
                    ├─► POST  → h.PostEvents
                    └─► GET   → h.GetEvents
                    └─► *     → 405
    │
    └─► return mux  ──► cmd/platform wraps in http.Server(:8081)
```

The store is created first and passed by interface into the handlers, so the handlers layer has no direct knowledge of the concrete `*store.Store` type.

## Request-Time: POST /events

An inbound POST request traverses from the network down to the store and back up as a response.

```
HTTP client
    │
    ▼
http.Server (:8081)
    │
    ▼
mux.HandleFunc("/events") — dispatch switch
    │  r.Method == POST
    ▼
handlers.PostEvents(w, r)
    │ decode JSON body → Event{Type, UserID, Payload}
    │ validate Type against allowlist {"click-mouse","submit-form"}
    │        └─► 400 Bad Request on invalid type or malformed JSON
    ▼
store.Storer.Add(event)
    │ acquire write lock
    │ generate UUID → event.ID
    │ append to []Event slice
    │ release write lock
    └─► return populated Event
    │
    ▼
handlers.PostEvents — JSON-encode response
    └─► 201 Created  {"id":…,"type":…,"userId":…,"payload":…}
```

## Request-Time: GET /events

```
HTTP client
    │
    ▼
http.Server (:8081)
    │
    ▼
mux.HandleFunc("/events") — dispatch switch
    │  r.Method == GET
    ▼
handlers.GetEvents(w, r)
    ▼
store.Storer.List()
    │ acquire read lock
    │ copy []Event slice
    │ release read lock
    └─► return []Event (may be empty)
    │
    ▼
handlers.GetEvents — JSON-encode response
    └─► 200 OK  [{…}, …]  ([] when empty)
```

## Cross-Component Coupling

The only coupling between `handlers` and `store` is the `store.Storer` interface:

```
handlers.Handler
    └─ store  store.Storer   ← interface boundary
                    ▲
              *store.Store   ← concrete type, injected at construction
```

`handlers` imports `store` for the interface definition only; `store` has no import of `handlers`. This one-way dependency keeps the sub-packages independently testable.
