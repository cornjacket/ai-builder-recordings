# store.go

Purpose: Implements the in-memory, concurrency-safe event store for the metrics service, defining the Event model, the Storer interface, and the Store concrete type.

Tags: data-access, model, interface

## Public API

### `type Event`

```go
type Event struct {
    ID      string
    Type    string
    UserID  string
    Payload map[string]interface{}
}
```

Domain type representing a single metrics event. `ID` is assigned by `Add` via UUID generation; callers supply `Type`, `UserID`, and `Payload`.

### `type Storer`

```go
type Storer interface {
    Add(e Event) Event
    List() []Event
}
```

Interface that the HTTP handler layer depends on. Allows test doubles to be injected in place of the concrete `Store`.

### `func New() *Store`

Returns a zero-value `Store` ready for use. No configuration is required.

### `func (s *Store) Add(e Event) Event`

Assigns a fresh UUID to the event's `ID` field, appends it to the in-memory slice under a write lock, and returns the stored copy with the assigned ID.

### `func (s *Store) List() []Event`

Returns a shallow copy of all stored events in insertion order, read under a read lock. The returned slice is independent of the internal slice — callers may mutate it safely.

## Key Internals

### Concurrency model

`Store` protects its `events` slice with a `sync.RWMutex`. `Add` acquires a full write lock; `List` acquires only a read lock, allowing concurrent reads. The shallow copy in `List` prevents callers from holding a reference into the internal slice after the lock is released.

### Compile-time interface assertion

```go
var _ Storer = (*Store)(nil)
```

Causes a compile error if `*Store` ever stops satisfying `Storer`, keeping the concrete type and interface in sync without runtime overhead.

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/google/uuid` | Generates universally unique string IDs for each stored event |
