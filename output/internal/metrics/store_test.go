// Package metrics provides an in-memory event recording service for frontend user interactions.
//
// Purpose: Unit tests for EventStore Add and List operations, including concurrent-safety checks.
// Tags: implementation, metrics
package metrics

import (
	"sync"
	"testing"
)

func TestAdd_AssignsID(t *testing.T) {
	s := NewEventStore()
	e := s.Add(Event{Type: "click-mouse", UserID: "u1", Payload: map[string]interface{}{}})
	if e.ID == "" {
		t.Fatal("expected non-empty ID after Add")
	}
}

func TestAdd_ReturnsStoredFields(t *testing.T) {
	s := NewEventStore()
	payload := map[string]interface{}{"key": "val"}
	e := s.Add(Event{Type: "submit-form", UserID: "u2", Payload: payload})
	if e.Type != "submit-form" {
		t.Errorf("got type %q, want %q", e.Type, "submit-form")
	}
	if e.UserID != "u2" {
		t.Errorf("got userId %q, want %q", e.UserID, "u2")
	}
}

func TestList_EmptyIsNotNil(t *testing.T) {
	s := NewEventStore()
	events := s.List()
	if events == nil {
		t.Fatal("List() returned nil, want empty slice")
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestList_ReturnsCopyInOrder(t *testing.T) {
	s := NewEventStore()
	s.Add(Event{Type: "click-mouse", UserID: "u1", Payload: map[string]interface{}{}})
	s.Add(Event{Type: "submit-form", UserID: "u2", Payload: map[string]interface{}{}})

	events := s.List()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Type != "click-mouse" {
		t.Errorf("first event type got %q, want %q", events[0].Type, "click-mouse")
	}
	if events[1].Type != "submit-form" {
		t.Errorf("second event type got %q, want %q", events[1].Type, "submit-form")
	}
}

func TestAdd_UniqueIDs(t *testing.T) {
	s := NewEventStore()
	e1 := s.Add(Event{Type: "click-mouse", UserID: "u1", Payload: map[string]interface{}{}})
	e2 := s.Add(Event{Type: "submit-form", UserID: "u2", Payload: map[string]interface{}{}})
	if e1.ID == e2.ID {
		t.Errorf("expected distinct IDs, both got %q", e1.ID)
	}
}

func TestAdd_Concurrent(t *testing.T) {
	s := NewEventStore()
	var wg sync.WaitGroup
	n := 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Add(Event{Type: "click-mouse", UserID: "u", Payload: map[string]interface{}{}})
		}()
	}
	wg.Wait()

	events := s.List()
	if len(events) != n {
		t.Errorf("expected %d events after concurrent adds, got %d", n, len(events))
	}
}
