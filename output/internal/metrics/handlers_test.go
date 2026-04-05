// Package metrics provides an in-memory event recording service for frontend user interactions.
//
// Purpose: Unit tests for PostEvent and GetEvents handlers using net/http/httptest.
// Tags: implementation, metrics
package metrics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostEvent_Success(t *testing.T) {
	store := NewEventStore()
	body := `{"type":"click-mouse","userId":"u1","payload":{}}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	w := httptest.NewRecorder()

	PostEvent(store)(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var e Event
	if err := json.NewDecoder(w.Body).Decode(&e); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if e.ID == "" {
		t.Error("expected non-empty id")
	}
	if e.Type != "click-mouse" {
		t.Errorf("got type %q, want click-mouse", e.Type)
	}
	if e.UserID != "u1" {
		t.Errorf("got userId %q, want u1", e.UserID)
	}
}

func TestPostEvent_InvalidType(t *testing.T) {
	store := NewEventStore()
	body := `{"type":"unknown-type","userId":"u1","payload":{}}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	w := httptest.NewRecorder()

	PostEvent(store)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPostEvent_MalformedJSON(t *testing.T) {
	store := NewEventStore()
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString("{bad json"))
	w := httptest.NewRecorder()

	PostEvent(store)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPostEvent_WrongMethod(t *testing.T) {
	store := NewEventStore()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	PostEvent(store)(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestGetEvents_EmptyReturnsArray(t *testing.T) {
	store := NewEventStore()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	GetEvents(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	// json.Encoder appends a newline; trim it
	var events []Event
	if err := json.Unmarshal([]byte(body), &events); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if events == nil {
		t.Error("expected non-nil slice, got null")
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestGetEvents_WrongMethod(t *testing.T) {
	store := NewEventStore()
	req := httptest.NewRequest(http.MethodPost, "/events", nil)
	w := httptest.NewRecorder()

	GetEvents(store)(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestGetEvents_ReturnsStoredEvents(t *testing.T) {
	store := NewEventStore()

	for _, body := range []string{
		`{"type":"click-mouse","userId":"u1","payload":{}}`,
		`{"type":"submit-form","userId":"u2","payload":{"key":"val"}}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		w := httptest.NewRecorder()
		PostEvent(store)(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("setup POST failed with %d", w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()
	GetEvents(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var events []Event
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Type != "click-mouse" {
		t.Errorf("first event type got %q, want click-mouse", events[0].Type)
	}
	if events[1].Type != "submit-form" {
		t.Errorf("second event type got %q, want submit-form", events[1].Type)
	}
	if events[0].ID == events[1].ID {
		t.Error("expected distinct IDs for two events")
	}
}
