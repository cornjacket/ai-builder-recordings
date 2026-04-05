# store

Purpose: In-memory, concurrency-safe event store for the metrics service — owns the Event model, the Storer interface, and UUID assignment on write.

Tags: data-access, model, interface

## Overview

The `store` package is the persistence layer of the metrics service. It holds all events in a `sync.RWMutex`-protected slice and assigns UUIDs at write time. The `Storer` interface decouples the HTTP handler layer from the concrete implementation, enabling test injection.

## File Index

| File | Description | Doc |
|------|-------------|-----|
| [store.go](store.go) | Event model, Storer interface, and Store implementation | [store.go.md](store.go.md) |

## Invariants

- **No persistence.** All events are lost when the process exits; there is no flush or recovery mechanism.
- **UUID assigned on Add.** Callers must not set `Event.ID` before calling `Add` — the store overwrites it unconditionally.
- **List returns a copy.** Mutating the slice returned by `List` does not affect the store's internal state.
- **Append-only.** There is no Delete or Update operation; the event log grows monotonically for the lifetime of the process.
