// Package metrics provides an in-memory event recording service for frontend user interactions.
//
// Purpose: Provides thread-safe storage and HTTP handlers for recording and listing user interaction events.
// Tags: implementation, metrics
package metrics

import (
	"sync"

	"github.com/google/uuid"
)

// Event represents a single frontend user interaction event.
type Event struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	UserID  string                 `json:"userId"`
	Payload map[string]interface{} `json:"payload"`
}

// EventStore is a thread-safe in-memory store for Event records.
type EventStore struct {
	mu     sync.RWMutex
	events []Event
}

// NewEventStore returns an initialised, empty EventStore.
func NewEventStore() *EventStore {
	return &EventStore{}
}

// Add assigns a UUID to the event, appends it to the store, and returns the stored event.
func (s *EventStore) Add(e Event) Event {
	e.ID = uuid.New().String()
	s.mu.Lock()
	s.events = append(s.events, e)
	s.mu.Unlock()
	return e
}

// List returns a copy of all stored events. The returned slice is never nil.
func (s *EventStore) List() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Event, len(s.events))
	copy(result, s.events)
	return result
}
