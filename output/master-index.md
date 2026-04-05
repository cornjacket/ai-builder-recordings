# output

## output

| File | Tags | Description |
|------|------|-------------|
| [data-flow.md](data-flow.md) | data-flow, architecture | Cross-component data-flow for the platform monolith — how cmd/platform constructs both services at startup and how HTTP requests are dispatched to each at runtime. |

## cmd

| File | Tags | Description |
|------|------|-------------|
| [README.md](cmd/README.md) | overview, architecture | Directory-level index for the cmd subtree — the runnable entry point of the platform monolith that wires together the metrics and IAM services. |
| [data-flow.md](cmd/data-flow.md) | data-flow, architecture | Describes how the platform binary wires together the internal service packages and routes incoming HTTP traffic to each. |

## cmd/platform

| File | Tags | Description |
|------|------|-------------|
| [README.md](cmd/platform/README.md) | main | Binary entry point that wires together the metrics and IAM internal packages and serves each on a dedicated port. |
| [main.go.md](cmd/platform/main.go.md) | main | Entry point for the platform binary — instantiates the metrics and IAM HTTP handlers and starts each on its own port concurrently. |

## internal

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/README.md) | overview, architecture | Directory index for the internal service packages — two independent HTTP services (iam and metrics) that are composed and served by the cmd/platform binary. |
| [data-flow.md](internal/data-flow.md) | data-flow, architecture | Describes how the iam and metrics services are constructed and how HTTP requests flow through each at runtime. |

## internal/iam

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/README.md) | architecture, overview | Identity and access management composite — exposes user lifecycle and RBAC authorisation as a single `http. |
| [data-flow.md](internal/iam/data-flow.md) | data-flow | Describes how data moves between the three components of the IAM composite — lifecycle, authz, and the iam. |
| [iam.go.md](internal/iam/iam.go.md) | architecture | Entry-point wiring file for the IAM composite — `New()` instantiates both sub-package handlers and registers all ten HTTP routes on a single `ServeMux`. |

### internal/iam/authz

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/authz/README.md) | api, model | RBAC sub-package — manages roles, user-role assignments, and permission checks via an in-memory store. |
| [authz.go.md](internal/iam/authz/authz.go.md) | api, model | Implements the RBAC sub-package — defines role and assignment models, an in-memory store protected by a read/write mutex, and an HTTP Handler that registers five permission-management routes. |

### internal/iam/lifecycle

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/lifecycle/README.md) | api, data-access, model | User lifecycle sub-package for the IAM service — owns user registration, lookup, deletion, |
| [lifecycle.go.md](internal/iam/lifecycle/lifecycle.go.md) | api, data-access, model | Implements user CRUD and session token management for the IAM lifecycle sub-package, |

## internal/metrics

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/metrics/README.md) | architecture, overview | HTTP ingestion service that records frontend user interaction events in memory and exposes them for retrieval. |
| [data-flow.md](internal/metrics/data-flow.md) | data-flow, architecture | Describes how data moves between the store and handlers sub-packages within the metrics composite, covering both startup wiring and request-time paths. |
| [metrics.go.md](internal/metrics/metrics.go.md) | architecture, data-flow | Documents the package-level New() constructor that wires the store and handlers sub-packages into a single http. |

### internal/metrics/handlers

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/metrics/handlers/README.md) | api | HTTP handler layer for the metrics service — decodes, validates, and routes POST /events and GET /events requests to the store sub-package. |
| [handlers.go.md](internal/metrics/handlers/handlers.go.md) | api | HTTP handlers for the metrics service, implementing POST /events (create) and GET /events (list) by delegating to a store. |

### internal/metrics/store

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/metrics/store/README.md) | data-access, model, interface | In-memory, concurrency-safe event store for the metrics service — owns the Event model, the Storer interface, and UUID assignment on write. |
| [store.go.md](internal/metrics/store/store.go.md) | data-access, model, interface | Implements the in-memory, concurrency-safe event store for the metrics service, defining the Event model, the Storer interface, and the Store concrete type. |

