# metrics

Purpose: Describes the design of the in-memory metrics package that records and lists frontend user interaction events.
Exposes two HTTP endpoints on port 8081, implemented as a reusable `http.Handler` returned by `routes.go`.

Tags: architecture, design

## File Index

| File | Description |
|------|-------------|
| `store.go` | `EventStore` — thread-safe in-memory slice of `Event` records with `Add` and `List` methods |
| `handlers.go` | `PostEvent` and `GetEvents` — `http.HandlerFunc` factories wired to an `EventStore` |
| `routes.go` | `NewRouter` — assembles and returns the `http.Handler` for the metrics service |
| `store_test.go` | Unit tests for `EventStore.Add` and `EventStore.List` |
| `handlers_test.go` | Unit tests for `PostEvent` and `GetEvents` using `net/http/httptest` |

## Overview

### Data model

```go
type Event struct {
    ID      string                 `json:"id"`
    Type    string                 `json:"type"`
    UserID  string                 `json:"userId"`
    Payload map[string]interface{} `json:"payload"`
}
```

Field names match the spec verbatim (`id`, `type`, `userId`, `payload`).

### EventStore (`store.go`)

```go
type EventStore struct {
    mu     sync.RWMutex
    events []Event
}

func NewEventStore() *EventStore
func (s *EventStore) Add(e Event) Event   // caller supplies all fields except ID; Add assigns UUID and appends
func (s *EventStore) List() []Event       // returns a shallow copy; never returns nil
```

`Add` holds a write lock; `List` holds a read lock. `List` returns `[]Event{}` (not nil) so `json.Marshal` produces `[]` rather than `null`.

### Handlers (`handlers.go`)

Both are constructor functions returning `http.HandlerFunc`:

- **`PostEvent(store *EventStore) http.HandlerFunc`**
  - Decodes JSON body into a request struct `{Type, UserID, Payload}`.
  - Validates `Type` ∈ `{"click-mouse", "submit-form"}` — returns `400 Bad Request` otherwise.
  - Calls `store.Add`, writes the returned `Event` as JSON with status `201 Created`.
  - Returns `400` on JSON decode failure; `405 Method Not Allowed` if not POST.

- **`GetEvents(store *EventStore) http.HandlerFunc`**
  - Calls `store.List`, writes the result as JSON with status `200 OK`.
  - Returns `405 Method Not Allowed` if not GET.

### Routes (`routes.go`)

```go
func NewRouter(store *EventStore) http.Handler
```

Registers `POST /events` and `GET /events` on a `http.NewServeMux()`. Method enforcement is handled inside each handler, not in the mux pattern.

### Dependencies

- `github.com/google/uuid` — UUID v4 generation for event IDs.
- Standard library only otherwise (`net/http`, `encoding/json`, `sync`).

### Testing strategy

- `store_test.go`: verifies `Add` assigns a non-empty ID, stores the event, and `List` returns it; verifies concurrent `Add` calls do not race (run with `-race`).
- `handlers_test.go`: uses `httptest.NewRecorder` and `httptest.NewRequest`; covers 201 on valid POST, 400 on invalid type, 400 on bad JSON, 200 array on GET (empty and non-empty), and 405 on wrong method for each endpoint.

## Documentation

### Implementation Notes
| File | Description |
|------|-------------|
| handlers.go.md | HTTP handler constructors for the metrics service, implementing POST and GET for the /events endpoint |
| store.go.md | Thread-safe in-memory store for frontend user interaction events, with Add and List operations |

