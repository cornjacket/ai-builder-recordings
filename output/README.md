# platform

Purpose: Repository-level index for the platform monolith — a single Go binary that serves two independent HTTP services (IAM on :8082, metrics on :8081) from one process.
The module `github.com/cornjacket/platform` (Go 1.22) has no external database; all runtime state is in-process memory.

Tags: overview, architecture

## Overview

The platform monolith is composed of a runnable entry-point layer (`cmd`) and two self-contained service packages (`internal/iam`, `internal/metrics`). The binary in `cmd/platform` is the only point of composition: it calls each service's `New()` constructor and mounts the returned `http.Handler` on a dedicated port via separate goroutines.

## Components

| Component | Responsibility | Docs |
|-----------|---------------|------|
| [`cmd/`](cmd/README.md) | Runnable binary that instantiates and serves both internal services | [README](cmd/README.md) |
| [`internal/`](internal/README.md) | Two independent HTTP service packages (iam and metrics) with no Go-level coupling to each other | [README](internal/README.md) |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how the binary wires the two services at startup and how requests are routed at runtime
