# iam

Purpose: Identity and access management composite — exposes user lifecycle and RBAC authorisation as a single `http.Handler` on port 8082.

Tags: architecture, overview

## Components

| Component | Description |
|-----------|-------------|
| [lifecycle/](lifecycle/README.md) | In-memory user CRUD, bcrypt password hashing, and bearer-token session management |
| [authz/](authz/README.md) | In-memory RBAC — role definitions, user-role assignments, and permission checks |
| [iam.go](iam.go.md) | Wiring layer — `New()` composes both sub-packages onto a single `ServeMux` |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how requests and data move between lifecycle, authz, and the wiring layer

## Overview

The `iam` package is a composite of two independent sub-packages, assembled by `iam.go`. `lifecycle` owns user registration, authentication, and session tokens. `authz` owns role definitions and permission resolution. Neither sub-package imports the other; cross-component coupling flows through the HTTP surface via shared string user IDs.

### Route Summary

| Method | Path | Handler |
|--------|------|---------|
| POST | /users | lifecycle |
| GET | /users/{id} | lifecycle |
| DELETE | /users/{id} | lifecycle |
| POST | /auth/login | lifecycle |
| POST | /auth/logout | lifecycle |
| POST | /roles | authz |
| GET | /roles | authz |
| POST | /users/{id}/roles | authz |
| GET | /users/{id}/roles | authz |
| POST | /authz/check | authz |

### Design Constraints

- No external database; all state lives in process memory and is lost on restart.
- Go module: `github.com/cornjacket/platform`; packages at `internal/iam`, `internal/iam/lifecycle`, `internal/iam/authz`.
- `iam.go` is the only place the two sub-packages are composed; neither sub-package imports the other.
