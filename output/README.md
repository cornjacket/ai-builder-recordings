# user-service

Purpose: Root directory index for the user service — a self-contained Go HTTP server exposing a JSON CRUD API for user records, backed by a thread-safe in-memory store.
All application logic lives in the `internal/` sub-tree; `main.go` is the sole entry point and wiring site.

Tags: architecture, overview

## Overview

The user service is a single-binary Go application. `main.go` wires the store and handler packages together and starts an HTTP server on `:8080`. All user-management logic is encapsulated within the `internal/` package boundary, preventing accidental import from outside the module.

## Components

| Component | Responsibility | Doc |
|-----------|----------------|-----|
| `main.go` | Composition root — constructs store and handler, registers routes, starts `:8080` | [main.go.md](main.go.md) |
| `internal/` | Go internal-package boundary enclosing the full user CRUD service sub-tree | [internal/README.md](internal/README.md) |

## Synthesis Docs

- [data-flow.md](data-flow.md) — startup sequence and end-to-end request path from `main.go` through the internal boundary to the HTTP response
