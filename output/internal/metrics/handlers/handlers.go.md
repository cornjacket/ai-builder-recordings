# handlers.go

Purpose: HTTP handlers for the metrics service, implementing POST /events (create) and GET /events (list) by delegating to a store.Storer dependency.
Accepts and returns JSON-encoded event objects; validates the event type against an allowlist of two known values.

Tags: api

## Public API

### `type Handler struct`

Holds a `store.Storer` interface value. All HTTP handler methods are defined on this type.

### `func New(s store.Storer) *Handler`

Constructs a `Handler` wired to the given `Storer`. This is the only way to obtain a `Handler`; the struct has no exported fields.

### `func (h *Handler) PostEvents(w http.ResponseWriter, r *http.Request)`

Handles `POST /events`. Decodes a JSON body into an `eventRequest`, validates that `type` is one of `"click-mouse"` or `"submit-form"`, stores the event via `h.store.Add`, and responds with `201 Created` and the stored event (including its generated ID) as JSON.

Returns `400 Bad Request` if the body is not valid JSON or if the `type` field is not in the allowlist.

### `func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request)`

Handles `GET /events`. Fetches all events from `h.store.List` and writes them as a JSON array with `200 OK`. Returns an empty array when no events exist.

## Key Internals

### `type eventRequest struct`

Unexported request envelope decoded from the POST body. Fields: `Type` (string), `UserID` (string), `Payload` (free-form `map[string]interface{}`).

### `type eventResponse struct`

Unexported response envelope used for both POST and GET responses. Mirrors `store.Event` but is kept separate from the store type to avoid coupling the HTTP shape directly to the storage model.

### Event-type allowlist in `PostEvents`

Only `"click-mouse"` and `"submit-form"` are accepted. Any other value results in `400 Bad Request`. This is an inline string comparison, not a lookup table.

## Dependencies

- `github.com/cornjacket/platform/internal/metrics/store` — provides the `Storer` interface and the `Event` type used for storage and retrieval.
