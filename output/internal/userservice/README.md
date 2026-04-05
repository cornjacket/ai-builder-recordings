# userservice

Purpose: Composite package providing a complete in-process user CRUD service — an HTTP handler layer backed by a thread-safe in-memory store.
The two sub-packages together form the full request-to-storage path for user records.

Tags: architecture, overview

## Overview

`userservice` contains the `handlers` and `store` sub-packages. Together they implement all user-management functionality: routing HTTP requests, performing in-memory data operations, and returning JSON responses. `main.go` at the project root wires these two packages into a running server on `:8080`.

## Components

| Sub-package | Responsibility | Doc |
|-------------|----------------|-----|
| [`handlers/`](handlers/README.md) | Routes `/users` HTTP requests to store operations and encodes JSON responses |
| [`store/`](store/README.md) | Thread-safe in-memory UUID-keyed CRUD store for `User` records |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how an HTTP request travels from the mux through handlers into the store and back
