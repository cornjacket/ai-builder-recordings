# store

## Goal

Thread-safe in-memory store for User records. User model: {"id": string (UUID, assigned on create), "name": string, "email": string}. Exposes: Create(name, email string) User; Get(id string) (User, bool); Update(id, name, email string) (User, bool); Delete(id string) bool. Uses sync.RWMutex internally. No external dependencies beyond stdlib.

## Context

### Level 1 — 5594e1-0000-build-1


## Design

**Language:** Go, stdlib only (`sync`, `crypto/rand`, `fmt`).

**Package path:** `internal/userservice/store`

**Files to produce:**
- `store.go` — all types and methods
- `store_test.go` — table-driven + concurrency tests

**User struct:**
```go
type User struct {
    ID    string
    Name  string
    Email string
}
```

**Store struct:**
```go
type Store struct {
    mu      sync.RWMutex
    records map[string]User
}

func New() *Store {
    return &Store{records: make(map[string]User)}
}
```

**UUID generation (no external deps):**
```go
func newUUID() string {
    var b [16]byte
    _, _ = io.ReadFull(rand.Reader, b[:])
    b[6] = (b[6] & 0x0f) | 0x40 // version 4
    b[8] = (b[8] & 0x3f) | 0x80 // variant 10
    return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
        b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
```

**Method contracts:**
- `Create`: Lock → generate UUID → store → Unlock → return User
- `Get`: RLock → lookup → RUnlock → return (User, found)
- `Update`: Lock → lookup; if missing return (User{}, false); overwrite → Unlock → return (updated User, true)
- `Delete`: Lock → lookup; if missing return false; delete → Unlock → return true

All locks released via `defer` immediately after acquisition.

**Constraints:**
- No third-party packages; only `sync`, `crypto/rand`, `io`, `fmt`.
- The map is always initialised by `New()` — no nil-map panics possible.

## Acceptance Criteria

1. `Create("Alice", "alice@example.com")` returns a `User` with non-empty `ID`, `Name == "Alice"`, `Email == "alice@example.com"`.
2. Two successive `Create` calls return `User` values with distinct `ID` fields.
3. `Get` on an ID returned by `Create` returns the same `User` and `true`.
4. `Get` on an unknown ID returns the zero `User` and `false`.
5. `Update` on an existing ID returns a `User` with the new name and email, and `true`; the updated values are visible on subsequent `Get`.
6. `Update` on an unknown ID returns the zero `User` and `false`; the store is unmodified.
7. `Delete` on an existing ID returns `true`; a subsequent `Get` for that ID returns `false`.
8. `Delete` on an unknown ID returns `false`.
9. Running the test binary with `-race` reports no data races when 50 goroutines concurrently call `Create`, `Get`, `Update`, and `Delete`.

## Test Command

```
cd /Users/david/Go/src/github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/output && go test -race ./internal/userservice/store/...
```

