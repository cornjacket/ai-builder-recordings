# Task: doc-platform-monolith

| Field       | Value                  |
|-------------|------------------------|
| Task-type   | USER-TASK              |
| Status      | in-progress             |
| Epic        | main               |
| Tags        | —               |
| Priority    | —           |
| Category    | —                      |
| Created     | 2026-04-05            |
| Completed   | —                      |
| Next-subtask-id | 0001 |

## Goal
Generate documentation for the platform-monolith source tree. Produce companion
`.md` files for every source file, a `README.md` in every directory, and
cross-component synthesis docs (`data-flow.md`, `api.md`) at composite levels.

The source tree is a Go HTTP service with two internal service packages:

```
cmd/platform/
  main.go                            ← service entry point (wires metrics + iam)
internal/
  metrics/
    metrics.go                       ← package wiring: store + handlers
    store/
      store.go                       ← in-memory event store
    handlers/
      handlers.go                    ← HTTP handlers for POST/GET /events
  iam/
    iam.go                           ← package wiring: lifecycle + authz
    lifecycle/
      lifecycle.go                   ← user CRUD + token auth
    authz/
      authz.go                       ← RBAC roles and permission checks
```

Some directories already have a `README.md` (`internal/iam/` and
`internal/metrics/`). The pipeline must co-exist with these pre-existing docs —
it should add companion `.md` files and sub-directory `README.md` files where
they are missing, and produce synthesis docs at composite levels.

## Context
This is a regression test for the ai-builder doc pipeline.

The pipeline must:
- Traverse the source tree recursively
- At each leaf directory (atomic): read source files, write companion `.md` files
  and a `README.md` describing the package
- At each composite directory (after all leaves complete): write cross-component
  synthesis docs (`data-flow.md`, `api.md`, `README.md`) using the completed
  sub-component handoff summaries
- Never modify source files — only create `.md` files
- Co-exist gracefully with pre-existing `.md` files in the source tree

The output directory IS the source tree. All generated `.md` files are written
inline alongside the source files they document.

## Subtasks

<!-- When a subtask is finished, run complete-task.sh --parent to mark it [x] before moving on. -->
<!-- subtask-list-start -->
- [x] [X-f8673c-0000-doc-1](X-f8673c-0000-doc-1/)
<!-- subtask-list-end -->

## Notes

_None._
