// Package handlers provides HTTP CRUD handlers for the user management API.
//
// Purpose: Tests for the HTTP handlers using an in-memory fake store and httptest.
// Tags: implementation, handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

// fakeStore is a minimal in-memory store for testing.
type fakeStore struct {
	users map[string]store.User
	seq   int
}

func newFakeStore() *fakeStore {
	return &fakeStore{users: make(map[string]store.User)}
}

func (f *fakeStore) Create(name, email string) store.User {
	f.seq++
	id := strings.Repeat("0", 7) + string(rune('0'+f.seq))
	u := store.User{ID: id, Name: name, Email: email}
	f.users[id] = u
	return u
}

func (f *fakeStore) Get(id string) (store.User, bool) {
	u, ok := f.users[id]
	return u, ok
}

func (f *fakeStore) Update(id, name, email string) (store.User, bool) {
	if _, ok := f.users[id]; !ok {
		return store.User{}, false
	}
	u := store.User{ID: id, Name: name, Email: email}
	f.users[id] = u
	return u, true
}

func (f *fakeStore) Delete(id string) bool {
	if _, ok := f.users[id]; !ok {
		return false
	}
	delete(f.users, id)
	return true
}

func TestCreateUser(t *testing.T) {
	h := New(newFakeStore())
	body := `{"name":"Alice","email":"alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
	var resp userResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty id")
	}
	if resp.Name != "Alice" || resp.Email != "alice@example.com" {
		t.Fatalf("unexpected body: %+v", resp)
	}
}

func TestCreateUserBadJSON(t *testing.T) {
	h := New(newFakeStore())
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{invalid"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetUser(t *testing.T) {
	fs := newFakeStore()
	h := New(fs)

	// Create a user first.
	u := fs.Create("Alice", "alice@example.com")

	req := httptest.NewRequest(http.MethodGet, "/users/"+u.ID, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp userResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != u.ID || resp.Name != "Alice" {
		t.Fatalf("unexpected body: %+v", resp)
	}
}

func TestGetUserNotFound(t *testing.T) {
	h := New(newFakeStore())
	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "{}" {
		t.Fatalf("expected body {}, got %s", w.Body.String())
	}
}

func TestUpdateUser(t *testing.T) {
	fs := newFakeStore()
	h := New(fs)
	u := fs.Create("Alice", "alice@example.com")

	body := `{"name":"Bob","email":"bob@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/users/"+u.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp userResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "Bob" || resp.Email != "bob@example.com" || resp.ID != u.ID {
		t.Fatalf("unexpected body: %+v", resp)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	h := New(newFakeStore())
	body := `{"name":"Bob","email":"bob@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/users/nonexistent", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUserBadJSON(t *testing.T) {
	h := New(newFakeStore())
	req := httptest.NewRequest(http.MethodPut, "/users/someid", strings.NewReader("{invalid"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	fs := newFakeStore()
	h := New(fs)
	u := fs.Create("Alice", "alice@example.com")

	req := httptest.NewRequest(http.MethodDelete, "/users/"+u.ID, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %s", w.Body.String())
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	h := New(newFakeStore())
	req := httptest.NewRequest(http.MethodDelete, "/users/nonexistent", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
