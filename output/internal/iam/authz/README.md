# authz

Purpose: RBAC sub-package — manages roles, user-role assignments, and permission checks via an in-memory store.

Tags: api, model

## Overview

The `authz` package provides a self-contained RBAC implementation. It exposes a single `Handler` that wires five HTTP routes for creating roles, assigning them to users, and evaluating whether a user holds a given permission string. All state is kept in an in-memory store guarded by a `sync.RWMutex`; there is no persistence layer.

## File Index

| File | Description | Doc |
|------|-------------|-----|
| `authz.go` | Role and assignment models, in-memory store, and HTTP Handler with five RBAC routes | [authz.go.md](authz.go.md) |

## Constraints and Invariants

- **No persistence** — all roles and assignments are lost on process restart.
- **Append-only assignments** — there is no API to revoke a role from a user once assigned.
- **nil-safe permissions** — `addRole` normalises a `nil` permissions slice to `[]string{}` so JSON responses never emit `null` for that field.
- **Role existence enforced at assignment** — `assignRole` rejects assignments that reference an unknown role ID, preventing dangling references in the assignments slice.
