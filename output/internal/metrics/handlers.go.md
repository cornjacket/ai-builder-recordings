Purpose: HTTP handler constructors for the metrics service, implementing POST and GET for the /events endpoint.
Each constructor closes over an *EventStore and returns an http.HandlerFunc.

Tags: implementation, metrics

## Key types

**validEventTypes** — package-level map used as a set for O(1) validation of the `type` field. Accepted values are `"click-mouse"` and `"submit-form"`.

## Main functions

**PostEvent(store) http.HandlerFunc** — decodes the JSON request body into an anonymous struct, rejects unknown event types with 400, calls `store.Add`, and writes the stored event as 201 JSON. Method is checked first; non-POST returns 405.

**GetEvents(store) http.HandlerFunc** — calls `store.List` and encodes the result as 200 JSON. Non-GET returns 405.

## Design decisions

- Method checks are handled inside each handler (in addition to the mux dispatch in routes.go) so that handlers remain independently testable without the router wrapping them.
- `json.NewEncoder(w).Encode(stored)` is used rather than `json.Marshal` + `w.Write` to avoid an intermediate buffer allocation.
- The anonymous decode struct mirrors only the fields the API accepts; clients cannot inject an ID through the request body.
