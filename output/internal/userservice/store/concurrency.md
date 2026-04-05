# concurrency

Purpose: Documents the locking strategy for the Store type, explaining which operations
acquire read vs. write locks and why.

Tags: architecture, design

## Locking Strategy

`Store` embeds a `sync.RWMutex`. The rule is simple: any operation that only reads the
map uses a read lock; any operation that modifies the map uses a full write lock.

| Method   | Lock type | Reason |
|----------|-----------|--------|
| `Get`    | `RLock` / `RUnlock` | Read-only; multiple goroutines may read concurrently |
| `Create` | `Lock` / `Unlock`   | Inserts a new key — map mutation |
| `Update` | `Lock` / `Unlock`   | Overwrites an existing key — map mutation |
| `Delete` | `Lock` / `Unlock`   | Removes a key — map mutation |

## Invariants

- The internal map is never `nil` after `New()` returns.
- `Create` always returns a `User` whose `ID` is a freshly generated UUID; it never
  returns a zero-value `User`.
- `Get`, `Update`, and `Delete` signal absence via the boolean second return value /
  boolean return value — they never panic on a missing key.
- Locks are always released via `defer` immediately after acquisition to prevent
  lock leaks on early returns.
