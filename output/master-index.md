# output

## internal/userservice

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/README.md) | architecture, design | Provides an in-memory HTTP user management service with CRUD operations over a REST API. |
| [theory-of-operation.md](internal/userservice/theory-of-operation.md) | architecture, design | Describes the data flow between the store and handlers components within the userservice package. |

### internal/userservice/handlers

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/handlers/README.md) | architecture, design | HTTP CRUD handlers for the user management API, wired to a Store interface. |
| [api.md](internal/userservice/handlers/api.md) | architecture, design | Full HTTP API contract for the user management handlers. |
| [handlers.go.md](internal/userservice/handlers/handlers.go.md) | implementation, handlers | HTTP CRUD handlers for the user management API, wired to a Store interface. |

### internal/userservice/store

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/userservice/store/README.md) | architecture, design | Defines the thread-safe in-memory store for User records used by the user-service. |
| [concurrency.md](internal/userservice/store/concurrency.md) | architecture, design | Documents the locking strategy for the Store type, explaining which operations |
| [store.md](internal/userservice/store/store.md) | implementation, store | Documents the design and key decisions of the thread-safe in-memory User store. |

