# metrics

## Goal

In-memory metrics service for recording and listing frontend user interaction events. Serves on port 8081.

HTTP API (verbatim from spec):
POST /events — body: {"type":"click-mouse"|"submit-form","userId":"<string>","payload":{}} → 201 with event object {"id":"<string>","type":"<string>","userId":"<string>","payload":{}};
GET /events → 200 JSON array of event objects, each with fields: id, type, userId, payload.

Data model: Event{id string, type string, userId string, payload map[string]interface{}}.
Store: thread-safe in-memory slice; IDs are generated UUIDs on POST.
Package layout: store.go (EventStore with Add/List), handlers.go (PostEvent, GetEvents), routes.go (returns http.Handler).
Unit tests must cover store and handler logic; handler tests use net/http/httptest.

## Context

### Level 1 — fbc6d9-0000-build-1


## Design

**Language:** Go

**Package:** `metrics` (import path `internal/metrics` within the host module)

**Files to produce:**

| File | Purpose |
|------|---------|
| `store.go` | `Event` struct, `EventStore` struct, `NewEventStore`, `Add`, `List` |
| `handlers.go` | `PostEvent(store) http.HandlerFunc`, `GetEvents(store) http.HandlerFunc` |
| `routes.go` | `NewRouter(store *EventStore) http.Handler` |
| `store_test.go` | Unit tests for store logic |
| `handlers_test.go` | Unit tests for handlers via httptest |

**Data model:**

```go
type Event struct {
    ID      string                 `json:"id"`
    Type    string                 `json:"type"`
    UserID  string                 `json:"userId"`
    Payload map[string]interface{} `json:"payload"`
}
```

**EventStore:**

```go
type EventStore struct {
    mu     sync.RWMutex
    events []Event
}
func NewEventStore() *EventStore
func (s *EventStore) Add(e Event) Event   // assigns uuid, appends, returns stored event
func (s *EventStore) List() []Event       // read lock; returns copy, never nil
```

**PostEvent handler:**
1. Check method == POST → 405 otherwise.
2. Decode body into `struct{ Type, UserID string; Payload map[string]interface{} }` → 400 on error.
3. Validate Type ∈ {"click-mouse", "submit-form"} → 400 with message otherwise.
4. Call store.Add with decoded fields; write 201 + JSON event.

**GetEvents handler:**
1. Check method == GET → 405 otherwise.
2. Call store.List; write 200 + JSON array.

**NewRouter:**
```go
func NewRouter(store *EventStore) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            PostEvent(store)(w, r)
        case http.MethodGet:
            GetEvents(store)(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })
    return mux
}
```

(Alternative: register method-dispatching directly in the mux pattern — either approach is acceptable.)

**Dependencies:**
- `github.com/google/uuid` for UUID v4 generation.
- Standard library: `net/http`, `encoding/json`, `sync`.

**Constraints:**
- `List` must initialise the return slice to `[]Event{}` (not nil) so JSON encodes as `[]` not `null`.
- No persistence; store is purely in-memory and resets on restart.
- The `NewRouter` return value is consumed by the integrate-platform wiring step which starts an `http.Server` on port 8081.

## Acceptance Criteria

1. `POST /events` with body `{"type":"click-mouse","userId":"u1","payload":{}}` returns HTTP 201 and a JSON object containing non-empty string `id`, `type` == `"click-mouse"`, `userId` == `"u1"`, and `payload` == `{}`.
2. `POST /events` with body `{"type":"submit-form","userId":"u2","payload":{"key":"val"}}` returns HTTP 201 and a JSON object with matching fields and a non-empty `id` distinct from the first event's id.
3. `POST /events` with `"type":"unknown-type"` returns HTTP 400.
4. `POST /events` with malformed JSON body returns HTTP 400.
5. `GET /events` after zero posts returns HTTP 200 and body `[]` (empty JSON array, not `null`).
6. `GET /events` after two posts returns HTTP 200 and a JSON array containing both event objects in insertion order, each with fields `id`, `type`, `userId`, `payload`.
7. `GET /events` with method POST on the GET-only path (or vice-versa where method dispatch is in the mux) returns HTTP 405.
8. Running `go test -race ./internal/metrics/...` from the module root passes with no failures and no race conditions detected.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/platform-monolith/output && go test -race ./internal/metrics/...
```

