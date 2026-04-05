# cmd

Purpose: Directory-level index for the cmd subtree — the runnable entry point of the platform monolith that wires together the metrics and IAM services.
The single binary under this directory starts both HTTP servers in the same process and is the sole consumer of the internal packages.

Tags: overview, architecture

## Components

| Component | Description | Docs |
|-----------|-------------|------|
| [`platform/`](platform/README.md) | Binary entry point that instantiates the metrics and IAM handlers and serves each on a dedicated port | [README](platform/README.md) |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how startup wiring and request routing flow through the `cmd` layer
