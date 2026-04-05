# Task: doc-user-service

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
Generate documentation for the user-service source tree. Produce companion `.md`
files for every source file, a `README.md` in every directory, and cross-component
synthesis docs (`data-flow.md`) at composite levels.

The source tree is a Go HTTP service with the following structure:

```
main.go                                    ← service entry point
internal/userservice/
  store/
    store.go                               ← in-memory user store
  handlers/
    handlers.go                            ← HTTP CRUD handlers
```

## Context
This is a regression test for the ai-builder doc pipeline.

The pipeline must:
- Traverse the source tree recursively
- At each leaf directory (atomic): read source files, write companion `.md` files
  and a `README.md` describing the package
- At each composite directory (after all leaves complete): write cross-component
  synthesis docs (`data-flow.md`, `README.md`) using the completed sub-component
  handoff summaries
- Never modify source files — only create `.md` files

The output directory IS the source tree. All generated `.md` files are written
inline alongside the source files they document.

## Subtasks

<!-- When a subtask is finished, run complete-task.sh --parent to mark it [x] before moving on. -->
<!-- subtask-list-start -->
- [ ] [6e440f-0000-doc-1](6e440f-0000-doc-1/)
<!-- subtask-list-end -->

## Notes

_None._
