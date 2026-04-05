# store

Purpose: Defines the thread-safe in-memory store for User records used by the user-service.
Exposes Create, Get, Update, and Delete operations backed by a sync.RWMutex-protected map.

Tags: architecture, design

## Files

| File | Description |
|------|-------------|
| `store.go` | `User` struct, `Store` struct, constructor `New()`, and all four CRUD methods |
| `store_test.go` | Table-driven tests covering all methods, concurrency, and not-found cases |

## Overview

The store is the single source of truth for User records within the process. It holds a
`map[string]User` protected by a `sync.RWMutex` so reads can proceed concurrently while
writes are serialised.

**User model**

```go
type User struct {
    ID    string
    Name  string
    Email string
}
```

**Constructor**

```go
func New() *Store
```

Returns an initialised `*Store` with a non-nil internal map.

**Method signatures**

```go
func (s *Store) Create(name, email string) User
func (s *Store) Get(id string) (User, bool)
func (s *Store) Update(id, name, email string) (User, bool)
func (s *Store) Delete(id string) bool
```

**UUID generation** uses `crypto/rand` to fill 16 random bytes, then sets the version-4
bits (byte 6 high nibble = `0100`) and variant bits (byte 8 high two bits = `10`) before
formatting as the canonical `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` string. No external
packages are required.

**Locking strategy** — see [concurrency.md](concurrency.md) for details.

## Documentation

### Design
| File | Description |
|------|-------------|
| concurrency.md | Documents the locking strategy for the Store type, explaining which operations |

