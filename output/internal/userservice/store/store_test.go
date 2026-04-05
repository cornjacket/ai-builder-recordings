// Package store provides a thread-safe in-memory store for User records.
//
// Purpose: Tests for the Store type covering all CRUD operations and concurrent safety.
// Tags: implementation, store
package store

import (
	"sync"
	"testing"
)

func TestCreate(t *testing.T) {
	s := New()
	u := s.Create("Alice", "alice@example.com")

	if u.ID == "" {
		t.Error("Create: ID must be non-empty")
	}
	if u.Name != "Alice" {
		t.Errorf("Create: Name = %q, want %q", u.Name, "Alice")
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Create: Email = %q, want %q", u.Email, "alice@example.com")
	}
}

func TestCreateDistinctIDs(t *testing.T) {
	s := New()
	u1 := s.Create("Alice", "alice@example.com")
	u2 := s.Create("Bob", "bob@example.com")
	if u1.ID == u2.ID {
		t.Errorf("Create: expected distinct IDs, got %q for both", u1.ID)
	}
}

func TestGetHit(t *testing.T) {
	s := New()
	created := s.Create("Alice", "alice@example.com")
	got, ok := s.Get(created.ID)
	if !ok {
		t.Fatal("Get: expected true for existing ID")
	}
	if got != created {
		t.Errorf("Get: got %+v, want %+v", got, created)
	}
}

func TestGetMiss(t *testing.T) {
	s := New()
	got, ok := s.Get("nonexistent-id")
	if ok {
		t.Error("Get: expected false for unknown ID")
	}
	if got != (User{}) {
		t.Errorf("Get: expected zero User, got %+v", got)
	}
}

func TestUpdateHit(t *testing.T) {
	s := New()
	u := s.Create("Alice", "alice@example.com")

	updated, ok := s.Update(u.ID, "Alicia", "alicia@example.com")
	if !ok {
		t.Fatal("Update: expected true for existing ID")
	}
	if updated.ID != u.ID {
		t.Errorf("Update: ID changed from %q to %q", u.ID, updated.ID)
	}
	if updated.Name != "Alicia" {
		t.Errorf("Update: Name = %q, want %q", updated.Name, "Alicia")
	}
	if updated.Email != "alicia@example.com" {
		t.Errorf("Update: Email = %q, want %q", updated.Email, "alicia@example.com")
	}

	// Verify persistence via Get.
	got, ok := s.Get(u.ID)
	if !ok {
		t.Fatal("Get after Update: expected true")
	}
	if got != updated {
		t.Errorf("Get after Update: got %+v, want %+v", got, updated)
	}
}

func TestUpdateMiss(t *testing.T) {
	s := New()
	s.Create("Alice", "alice@example.com") // populate store with one record

	got, ok := s.Update("nonexistent-id", "X", "x@example.com")
	if ok {
		t.Error("Update: expected false for unknown ID")
	}
	if got != (User{}) {
		t.Errorf("Update: expected zero User, got %+v", got)
	}
}

func TestDeleteHit(t *testing.T) {
	s := New()
	u := s.Create("Alice", "alice@example.com")

	if !s.Delete(u.ID) {
		t.Fatal("Delete: expected true for existing ID")
	}
	_, ok := s.Get(u.ID)
	if ok {
		t.Error("Get after Delete: expected false")
	}
}

func TestDeleteMiss(t *testing.T) {
	s := New()
	if s.Delete("nonexistent-id") {
		t.Error("Delete: expected false for unknown ID")
	}
}

func TestConcurrency(t *testing.T) {
	s := New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			u := s.Create("concurrent", "concurrent@example.com")
			if got, ok := s.Get(u.ID); ok {
				s.Update(got.ID, "updated", "updated@example.com")
			}
			s.Delete(u.ID)
		}()
	}

	wg.Wait()
}
