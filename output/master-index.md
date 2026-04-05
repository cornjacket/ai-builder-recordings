# output

## output

| File | Tags | Description |
|------|------|-------------|
| [data-flow.md](data-flow.md) | data-flow, architecture | End-to-end description of how the user service starts up and processes HTTP requests — from `main. |
| [main.go.md](main.go.md) | architecture, overview | Entry point for the user service — constructs the store and handler, registers routes on a ServeMux, and starts the HTTP server on :8080. |

## internal

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/README.md) | architecture, overview | Top-level internal package index for the user service — contains the userservice sub-tree, which implements the full HTTP CRUD stack for user records. |
| [data-flow.md](internal/data-flow.md) | data-flow, architecture | Describes how an inbound HTTP request crosses the internal/ package boundary and is processed by the userservice sub-tree before a response is returned to the caller. |

## internal/userservice

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/README.md) | architecture, overview | Composite package providing a complete in-process user CRUD service — an HTTP handler layer backed by a thread-safe in-memory store. |
| [data-flow.md](internal/userservice/data-flow.md) | data-flow, architecture | Describes how an HTTP request moves through the userservice composite — from the ServeMux registered in main. |

### internal/userservice/handlers

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/handlers/README.md) | api, interface | HTTP CRUD handler layer for the user service — routes incoming requests to store operations and encodes JSON responses. |
| [handlers.go.md](internal/userservice/handlers/handlers.go.md) | api, interface | Implements HTTP CRUD handlers for the user service, wiring POST/GET/PUT/DELETE `/users` routes onto a `*http. |

### internal/userservice/store

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/store/README.md) | data-access, model | Thread-safe in-memory data access layer that manages the full lifecycle of User records via UUID-keyed CRUD operations. |
| [store.go.md](internal/userservice/store/store.go.md) | data-access, model | Implements a thread-safe in-memory CRUD store for User records, using a `sync. |

