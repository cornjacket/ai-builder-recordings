# metrics

Purpose: HTTP ingestion service that records frontend user interaction events in memory and exposes them for retrieval. Listens on port 8081.

Tags: architecture, overview

## Overview

The metrics composite provides a single `http.Handler` (returned by `New()`) that wires an in-memory event store and an HTTP handler layer together onto one `ServeMux`. The only routes exposed are `POST /events` and `GET /events`.

## Components

| Component | Responsibility |
|-----------|----------------|
| [`handlers/`](handlers/README.md) | HTTP layer — decodes requests, enforces the event-type allowlist, delegates storage to a `store.Storer` |
| [`store/`](store/README.md) | Data-access layer — owns the `Event` type, `Storer` interface, and thread-safe in-memory storage |
| [`metrics.go`](metrics.go.md) | Package-level wiring — constructs store and handlers, registers routes, returns the composed `http.Handler` |

## Synthesis Docs

- [data-flow.md](data-flow.md) — startup wiring diagram and request-time paths for POST and GET

## API Contract

```
POST /events
  Request:  {"type":"click-mouse"|"submit-form","userId":string,"payload":{}}
  Response: 201 {"id":string,"type":string,"userId":string,"payload":{}}
  Errors:   400 on invalid JSON or unknown type

GET /events
  Response: 200 [{"id":string,"type":string,"userId":string,"payload":{}}]
            200 [] when no events have been recorded
```

## Data Model

```go
type Event struct {
    ID      string                 `json:"id"`
    Type    string                 `json:"type"`
    UserID  string                 `json:"userId"`
    Payload map[string]interface{} `json:"payload"`
}
```
