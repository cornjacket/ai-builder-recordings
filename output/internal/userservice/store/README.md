# store

Purpose: Thread-safe in-memory data access layer that manages the full lifecycle of User records via UUID-keyed CRUD operations.

Tags: data-access, model

## Overview

The `store` package owns all persistence logic for the user service. It exposes a concrete `*Store` type that the `handlers` package depends on through the `Store` interface defined in that package. All methods are safe for concurrent use.

## File Index

| File | Description | Doc |
|------|-------------|-----|
| `store.go` | `Store` struct, `User` model, and CRUD methods backed by a mutex-guarded map | [store.go.md](store.go.md) |

## Constraints and Invariants

- **ID assignment is internal** — callers must never set `User.ID` before calling `Create`; the store silently overwrites it with a fresh UUID.
- **Update preserves the original ID** — even if the caller supplies a different `ID` in the `user` argument to `Update`, the stored record retains the `id` parameter value.
- **No persistence** — all data lives in process memory and is lost on restart.
- **RWMutex strategy** — `Get` uses a read lock (concurrent reads are safe); `Create`, `Update`, and `Delete` use an exclusive write lock.
