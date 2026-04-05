Purpose: Thread-safe in-memory store for frontend user interaction events, with Add and List operations.
Designed to be embedded in an HTTP service that records and retrieves events concurrently.

Tags: implementation, store

## Key types

**Event** — the data model for a single user interaction. Fields `ID`, `Type`, `UserID`, and `Payload` map directly to the JSON API contract via struct tags.

**EventStore** — holds a `sync.RWMutex` and a `[]Event` slice. The mutex gates all reads and writes.

## Main functions

**NewEventStore** — returns a zero-value `*EventStore`. The slice starts nil but is initialised lazily by `append`.

**Add** — takes a write lock, generates a UUID v4 via `github.com/google/uuid`, appends the event, and returns the stored copy. UUID generation happens outside the lock — only the append is protected — keeping the critical section minimal.

**List** — takes a read lock and copies the slice via `make`+`copy`, guaranteeing callers receive a snapshot that cannot be modified by concurrent writes. Returns `[]Event{}` (never nil) so JSON encoding produces `[]` rather than `null`.

## Design decisions

- `sync.RWMutex` chosen over `sync.Mutex` because reads are expected to dominate; concurrent `GET /events` calls share the read lock without blocking each other.
- The return slice is initialised with `make([]Event, len(s.events))` before `copy` rather than `append(nil, s.events...)` to make the non-nil guarantee explicit and the intent clear.
