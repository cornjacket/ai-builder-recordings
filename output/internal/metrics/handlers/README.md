# handlers

Purpose: HTTP handler layer for the metrics service — decodes, validates, and routes POST /events and GET /events requests to the store sub-package.

Tags: api

## Overview

The `handlers` package owns the HTTP surface of the metrics service. It exposes two endpoints via a `Handler` struct that accepts a `store.Storer` interface, keeping the HTTP layer decoupled from the storage implementation.

## File Index

| File | Description | Docs |
|------|-------------|------|
| `handlers.go` | `Handler` type, `New` constructor, and `PostEvents`/`GetEvents` methods | [handlers.go.md](handlers.go.md) |

## Constraints and Invariants

- Event type is validated against a two-value allowlist (`"click-mouse"`, `"submit-form"`); all other values are rejected with `400 Bad Request`.
- The package depends on `store.Storer` by interface — any conforming store implementation may be injected.
- No authentication or authorisation is performed at this layer; that is the responsibility of the caller that wires the handler into a mux.
- `json.NewEncoder(w).Encode(...)` errors are not checked; network write failures are silent.
