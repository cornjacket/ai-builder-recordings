# iam

Purpose: Composite IAM service providing user lifecycle management, token-based authentication, and role-based access control over a single HTTP listener on port 8082.

Tags: architecture, design

The `iam` package wires together two self-contained sub-packages — `lifecycle` and `authz` — into one `http.ServeMux` that handles all ten IAM endpoints.

## Components

| Component | Description |
|-----------|-------------|
| [lifecycle](lifecycle/README.md) | In-memory user store and token-based authentication (POST /users, GET /users/{id}, DELETE /users/{id}, POST /auth/login, POST /auth/logout) |
| [authz](authz/README.md) | In-memory role store and RBAC permission checks (POST /roles, GET /roles, POST /users/{id}/roles, GET /users/{id}/roles, POST /authz/check) |

See [theory-of-operation.md](theory-of-operation.md) for data-flow and component interaction.
