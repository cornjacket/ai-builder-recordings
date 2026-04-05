# handlers

Purpose: HTTP CRUD handler layer for the user service — routes incoming requests to store operations and encodes JSON responses.

Tags: api, interface

## Overview

The `handlers` package wires four REST endpoints (`POST /users`, `GET /users/{id}`, `PUT /users/{id}`, `DELETE /users/{id}`) onto a `*http.ServeMux`. It depends on the `Store` interface rather than the concrete store type, keeping the two packages loosely coupled and independently testable.

## File Index

| File | Description | Doc |
|------|-------------|-----|
| `handlers.go` | `Handler` struct, `Store` interface, route registration, and per-verb handler functions | [handlers.go.md](handlers.go.md) |

## Constraints and Notes

- **Interface boundary:** `Handler` accepts any `Store` implementation; the production wiring in `main.go` passes `*store.Store`, but tests may supply a lightweight stub without importing the store package.
- **Go 1.22 routing:** `RegisterRoutes` relies on method-prefixed patterns (e.g. `"POST /users"`) introduced in Go 1.22. The service will not compile against earlier toolchains.
- **No middleware:** Request logging, authentication, and recovery are not handled here; they must be layered at the `ServeMux` level by the caller (`main.go`).
