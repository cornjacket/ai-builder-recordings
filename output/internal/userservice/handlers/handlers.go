// Package handlers provides HTTP CRUD handlers for the user management API.
//
// Purpose: Registers routes for POST/GET/PUT/DELETE /users endpoints and wires them to a Store.
// Responses are JSON-encoded; 400 on malformed input, 404 on missing records.
// Tags: implementation, handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

// Store is the interface the handlers depend on.
type Store interface {
	Create(name, email string) store.User
	Get(id string) (store.User, bool)
	Update(id, name, email string) (store.User, bool)
	Delete(id string) bool
}

// New registers all routes and returns the mux.
func New(s Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", handleCreate(s))
	mux.HandleFunc("GET /users/{id}", handleGet(s))
	mux.HandleFunc("PUT /users/{id}", handleUpdate(s))
	mux.HandleFunc("DELETE /users/{id}", handleDelete(s))
	return mux
}

type createRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handleCreate(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u := s.Create(req.Name, req.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userResponse{ID: u.ID, Name: u.Name, Email: u.Email})
	}
}

func handleGet(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		u, ok := s.Get(id)
		if !ok {
			writeNotFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userResponse{ID: u.ID, Name: u.Name, Email: u.Email})
	}
}

func handleUpdate(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req updateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u, ok := s.Update(id, req.Name, req.Email)
		if !ok {
			writeNotFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userResponse{ID: u.ID, Name: u.Name, Email: u.Email})
	}
}

func handleDelete(s Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !s.Delete(id) {
			writeNotFound(w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{}"))
}
