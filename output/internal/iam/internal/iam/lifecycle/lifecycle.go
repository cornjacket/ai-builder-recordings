// Package lifecycle provides in-memory user CRUD and token-based authentication for the IAM listener.
//
// Purpose: Manages user registration, lookup, deletion, and session token lifecycle via five HTTP handlers.
// Exposes Handler() to register all routes on a *http.ServeMux.
// Tags: implementation, lifecycle
package lifecycle

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User holds a single registered user's data.
type User struct {
	ID           string
	Username     string
	PasswordHash string
}

// UserStore is a concurrent-safe dual-indexed map of users.
type UserStore struct {
	mu         sync.RWMutex
	byID       map[string]*User
	byUsername map[string]*User
}

// TokenStore is a concurrent-safe map from token string to userID.
type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]string
}

// handler holds references to both stores and satisfies the routing methods.
type handler struct {
	users  *UserStore
	tokens *TokenStore
}

// Handler constructs both stores and returns a *http.ServeMux with all five routes registered.
func Handler() *http.ServeMux {
	h := &handler{
		users: &UserStore{
			byID:       make(map[string]*User),
			byUsername: make(map[string]*User),
		},
		tokens: &TokenStore{
			tokens: make(map[string]string),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.usersRoot)
	mux.HandleFunc("/users/", h.usersID)
	mux.HandleFunc("/auth/login", h.login)
	mux.HandleFunc("/auth/logout", h.logout)
	return mux
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// writeError writes a {"error":"<msg>"} JSON response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// usersRoot handles POST /users.
func (h *handler) usersRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to hash password")
		return
	}

	h.users.mu.Lock()
	defer h.users.mu.Unlock()

	if _, exists := h.users.byUsername[req.Username]; exists {
		writeError(w, http.StatusBadRequest, "username already taken")
		return
	}

	u := &User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		PasswordHash: string(hash),
	}
	h.users.byID[u.ID] = u
	h.users.byUsername[u.Username] = u

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":       u.ID,
		"username": u.Username,
	})
}

// usersID handles GET /users/{id} and DELETE /users/{id}.
func (h *handler) usersID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.users.mu.RLock()
		u, ok := h.users.byID[id]
		h.users.mu.RUnlock()

		if !ok {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"id":       u.ID,
			"username": u.Username,
		})

	case http.MethodDelete:
		h.users.mu.Lock()
		u, ok := h.users.byID[id]
		if !ok {
			h.users.mu.Unlock()
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		delete(h.users.byID, u.ID)
		delete(h.users.byUsername, u.Username)
		h.users.mu.Unlock()

		// Cascade: invalidate any tokens belonging to this user.
		h.tokens.mu.Lock()
		for tok, uid := range h.tokens.tokens {
			if uid == id {
				delete(h.tokens.tokens, tok)
			}
		}
		h.tokens.mu.Unlock()

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// login handles POST /auth/login.
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid request body")
		return
	}

	h.users.mu.RLock()
	u, ok := h.users.byUsername[req.Username]
	h.users.mu.RUnlock()

	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token := uuid.NewString()

	h.tokens.mu.Lock()
	h.tokens.tokens[token] = u.ID
	h.tokens.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// logout handles POST /auth/logout.
func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	h.tokens.mu.Lock()
	defer h.tokens.mu.Unlock()

	if _, ok := h.tokens.tokens[token]; !ok {
		writeError(w, http.StatusUnauthorized, "token not found")
		return
	}
	delete(h.tokens.tokens, token)

	w.WriteHeader(http.StatusOK)
}
