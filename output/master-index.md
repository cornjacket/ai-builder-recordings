# output

## internal/iam

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/README.md) | architecture, design | Composite IAM service providing user lifecycle management, token-based authentication, and role-based access control over a single HTTP listener on port 8082. |
| [iam.go.md](internal/iam/iam.go.md) | implementation, iam | Wires lifecycle and authz sub-muxes into the single IAM *http. |
| [theory-of-operation.md](internal/iam/theory-of-operation.md) | architecture, design | Describes the data-flow and component interactions within the IAM composite, showing how lifecycle and authz sub-packages collaborate to serve all ten IAM endpoints. |

#### internal/iam/internal/iam/authz

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/internal/iam/authz/README.md) | architecture, design | Provides in-memory role management and RBAC permission checks for the IAM listener. |
| [api.md](internal/iam/internal/iam/authz/api.md) | architecture, design | Defines the HTTP request and response contract for all five authz endpoints. |
| [authz.go.md](internal/iam/internal/iam/authz/authz.go.md) | implementation, authz | Documents the design and key decisions of the authz package's role store and RBAC handler. |

#### internal/iam/internal/iam/lifecycle

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/iam/internal/iam/lifecycle/README.md) | architecture, design | Provides in-memory user CRUD and token-based authentication for the IAM listener. |
| [api.md](internal/iam/internal/iam/lifecycle/api.md) | architecture, design | Defines the HTTP request and response contract for all five lifecycle endpoints. |
| [lifecycle.go.md](internal/iam/internal/iam/lifecycle/lifecycle.go.md) | implementation, lifecycle | Documents the design and implementation decisions for the lifecycle package, which owns user CRUD and session-token management for the IAM listener. |

## internal/metrics

| File | Tags | Description |
|------|------|-------------|
| [README.md](internal/metrics/README.md) | architecture, design | Describes the design of the in-memory metrics package that records and lists frontend user interaction events. |
| [handlers.go.md](internal/metrics/handlers.go.md) | implementation, metrics | HTTP handler constructors for the metrics service, implementing POST and GET for the /events endpoint. |
| [store.go.md](internal/metrics/store.go.md) | implementation, store | Thread-safe in-memory store for frontend user interaction events, with Add and List operations. |

