Purpose: Documents the design and key decisions of the thread-safe in-memory User store.
Covers types, methods, UUID generation, and locking strategy.

Tags: implementation, store

## Key Types

### `User`
Plain value struct with `ID`, `Name`, and `Email` string fields. Passed and returned by value throughout — no pointer aliasing issues.

### `Store`
Holds a `map[string]User` protected by a `sync.RWMutex`. Always created via `New()`, which initialises the map, so nil-map panics cannot occur.

## Methods

| Method | Lock | Behaviour |
|---|---|---|
| `Create(name, email)` | Write lock | Generates UUID v4, inserts into map, returns new User |
| `Get(id)` | Read lock | Returns (User, true) on hit; (User{}, false) on miss |
| `Update(id, name, email)` | Write lock | Returns (User{}, false) if ID unknown; overwrites entry and returns (User, true) on hit |
| `Delete(id)` | Write lock | Returns false if ID unknown; removes entry and returns true on hit |

All locks are acquired then immediately deferred for release — no manual `Unlock` calls.

## UUID Generation (`newUUID`)

Reads 16 random bytes from `crypto/rand`. Sets byte 6's high nibble to `0x4` (version 4) and byte 8's high two bits to `0b10` (RFC 4122 variant). Formats as `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` using `fmt.Sprintf` with byte-slice arguments.

## Design Decisions

- **Write lock for Update** — the check-then-write sequence must be atomic. A read lock for the lookup followed by upgrading to a write lock is not possible with `sync.RWMutex`, so a single write lock covers both steps.
- **Value semantics on User** — copying a small three-field struct is cheaper and safer than managing pointer lifetimes across concurrent callers.
- **`defer` on every lock acquisition** — eliminates the risk of forgetting to unlock on early returns (e.g. miss paths that return before the end of the function body).
