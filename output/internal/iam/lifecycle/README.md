# lifecycle

Purpose: User lifecycle sub-package for the IAM service — owns user registration, lookup, deletion,
and bearer-token session management entirely in memory.

Tags: api, data-access, model

## Overview

The `lifecycle` package exposes four HTTP endpoints (create/get/delete user; login/logout) through a single `Handler` backed by an in-memory `Store`. Passwords are hashed with bcrypt at default cost. Session tokens are opaque UUIDs with no expiry.

## File Index

| File | Description | Doc |
|---|---|---|
| `lifecycle.go` | Handler, Store, domain types, and all HTTP endpoint logic | [lifecycle.go.md](lifecycle.go.md) |

## Constraints and Invariants

- **No persistence.** All user and token state is lost on process restart.
- **Username uniqueness is enforced.** `addUser` rejects duplicate usernames with HTTP 409.
- **Tokens do not expire.** A token remains valid until an explicit `POST /auth/logout` call invalidates it; there is no TTL or background eviction.
- **Password hashes are never returned.** All JSON responses include only `id` and `username`; `PasswordHash` is never serialised to the wire.
- **Concurrency safety.** The `Store` uses a single `sync.RWMutex`; read-only operations take a read lock and mutating operations take a write lock, making all endpoints safe for concurrent use.
