# internal

Purpose: Directory index for the internal service packages — two independent HTTP services (iam and metrics) that are composed and served by the cmd/platform binary.
Each service exposes its own http.Handler and binds to a distinct port at runtime.

Tags: overview, architecture

## Overview

`internal` contains two self-contained service packages. `iam` handles identity and access management on port :8082; `metrics` handles event ingestion and retrieval on port :8081. Neither package imports the other — they are decoupled at the Go level and coordinated only by the binary in `cmd/platform`.

## Components

| Component | Responsibility |
|-----------|---------------|
| [iam/](iam/README.md) | Identity and access management: user lifecycle (registration, login, logout) and RBAC (roles, assignments, permission checks) composed onto a single ServeMux at :8082 |
| [metrics/](metrics/README.md) | Event collection and retrieval: in-memory event store with type-allowlist validation, exposed as POST /events and GET /events on a single ServeMux at :8081 |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how iam and metrics are constructed and how requests flow through each service at runtime
