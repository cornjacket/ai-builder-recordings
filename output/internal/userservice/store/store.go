// Package store provides a thread-safe in-memory store for User records.
//
// Purpose: Manages CRUD operations on User records using a sync.RWMutex-protected map.
// UUID v4 identifiers are generated via crypto/rand on each Create call.
// Tags: implementation, store
package store

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"
)

// User represents a user record.
type User struct {
	ID    string
	Name  string
	Email string
}

// Store is a thread-safe in-memory store for User records.
type Store struct {
	mu      sync.RWMutex
	records map[string]User
}

// New returns an initialised Store ready for use.
func New() *Store {
	return &Store{records: make(map[string]User)}
}

// Create inserts a new User with the given name and email, assigning a UUID v4 ID.
func (s *Store) Create(name, email string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{
		ID:    newUUID(),
		Name:  name,
		Email: email,
	}
	s.records[u.ID] = u
	return u
}

// Get returns the User with the given ID and true, or a zero User and false if not found.
func (s *Store) Get(id string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.records[id]
	return u, ok
}

// Update replaces the name and email of the User with the given ID.
// Returns the updated User and true on success, or a zero User and false if the ID is unknown.
func (s *Store) Update(id, name, email string) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[id]; !ok {
		return User{}, false
	}
	u := User{ID: id, Name: name, Email: email}
	s.records[id] = u
	return u, true
}

// Delete removes the User with the given ID, returning true on success or false if not found.
func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[id]; !ok {
		return false
	}
	delete(s.records, id)
	return true
}

// newUUID generates a random UUID v4.
func newUUID() string {
	var b [16]byte
	_, _ = io.ReadFull(rand.Reader, b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
