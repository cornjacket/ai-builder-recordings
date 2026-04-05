# store.go

Purpose: Implements a thread-safe in-memory CRUD store for User records, using a `sync.RWMutex`-guarded map and UUID generation for new record IDs.

Tags: data-access, model

## Public API

### Types

#### `User`
```go
type User struct {
    ID    string
    Name  string
    Email string
}
```
Domain type representing a single user record. `ID` is always a UUID string assigned by the store on creation.

#### `Store`
```go
type Store struct { ... }
```
Thread-safe in-memory store backed by a `map[string]User`. All exported methods are safe for concurrent use.

### Functions

#### `New() *Store`
Returns an initialised `*Store` with an empty internal map, ready for use.

#### `(*Store).Create(user User) User`
Ignores any `ID` on the input, assigns a fresh UUID, persists the record, and returns the populated `User`. Acquires a write lock for the duration.

#### `(*Store).Get(id string) (User, bool)`
Returns the stored `User` and `true` for the given `id`, or a zero `User` and `false` if the record does not exist. Acquires a read lock, allowing concurrent reads.

#### `(*Store).Update(id string, user User) (User, bool)`
Replaces the stored record for `id` with the supplied `user`, preserving the original `ID` field. Returns the updated `User` and `true` on success, or `(User{}, false)` if `id` is not found. Acquires a write lock.

#### `(*Store).Delete(id string) bool`
Removes the record for `id` and returns `true`, or `false` if the record does not exist. Acquires a write lock.

## Key Internals

`mu sync.RWMutex` guards `data`. Read-only methods (`Get`) use `RLock`/`RUnlock` to allow concurrent reads; mutating methods (`Create`, `Update`, `Delete`) use `Lock`/`Unlock`. UUID strings are generated via `uuid.NewString()` from the `github.com/google/uuid` package — callers have no influence over ID values.

## Dependencies

- `github.com/google/uuid` — UUID v4 generation for new record IDs
