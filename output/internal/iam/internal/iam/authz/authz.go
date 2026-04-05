// Package authz provides in-memory RBAC role management and permission-check handling.
//
// Purpose: Manages roles and user-role assignments with a permission-check endpoint.
// Exposes a *http.ServeMux for registration by the IAM listener.
// Tags: implementation, authz
package authz

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Role represents an access-control role with a set of named permissions.
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type handler struct {
	mu        sync.RWMutex
	roleStore map[string]Role     // keyed by role ID
	userRoles map[string][]string // keyed by user ID → []roleID
}

// Handler returns a *http.ServeMux with the five authz endpoints registered.
func Handler() *http.ServeMux {
	h := &handler{
		roleStore: make(map[string]Role),
		userRoles: make(map[string][]string),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/roles", h.handleRoles)
	mux.HandleFunc("/users/", h.handleUsers)
	mux.HandleFunc("/authz/check", h.handleAuthzCheck)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// handleRoles dispatches POST /roles and GET /roles.
func (h *handler) handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createRole(w, r)
	case http.MethodGet:
		h.listRoles(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *handler) createRole(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	perms := body.Permissions
	if perms == nil {
		perms = []string{}
	}
	role := Role{
		ID:          uuid.New().String(),
		Name:        body.Name,
		Permissions: perms,
	}
	h.mu.Lock()
	h.roleStore[role.ID] = role
	h.mu.Unlock()
	writeJSON(w, http.StatusCreated, role)
}

func (h *handler) listRoles(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	roles := make([]Role, 0, len(h.roleStore))
	for _, role := range h.roleStore {
		roles = append(roles, role)
	}
	h.mu.RUnlock()
	writeJSON(w, http.StatusOK, roles)
}

// handleUsers dispatches POST /users/{id}/roles and GET /users/{id}/roles.
func (h *handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Path: /users/{id}/roles
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "users" || parts[2] != "roles" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	userID := parts[1]
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user id required")
		return
	}
	switch r.Method {
	case http.MethodPost:
		h.assignRole(w, r, userID)
	case http.MethodGet:
		h.getUserRoles(w, r, userID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *handler) assignRole(w http.ResponseWriter, r *http.Request, userID string) {
	var body struct {
		RoleID string `json:"roleId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.RoleID == "" {
		writeError(w, http.StatusBadRequest, "roleId is required")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.roleStore[body.RoleID]; !ok {
		writeError(w, http.StatusNotFound, "role not found")
		return
	}
	for _, id := range h.userRoles[userID] {
		if id == body.RoleID {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	h.userRoles[userID] = append(h.userRoles[userID], body.RoleID)
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) getUserRoles(w http.ResponseWriter, r *http.Request, userID string) {
	h.mu.RLock()
	roleIDs := h.userRoles[userID]
	roles := make([]Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		if role, ok := h.roleStore[id]; ok {
			roles = append(roles, role)
		}
	}
	h.mu.RUnlock()
	writeJSON(w, http.StatusOK, roles)
}

// handleAuthzCheck handles POST /authz/check.
func (h *handler) handleAuthzCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		UserID     string `json:"userId"`
		Permission string `json:"permission"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	h.mu.RLock()
	roleIDs := h.userRoles[body.UserID]
	allowed := false
outer:
	for _, id := range roleIDs {
		role, ok := h.roleStore[id]
		if !ok {
			continue
		}
		for _, p := range role.Permissions {
			if p == body.Permission {
				allowed = true
				break outer
			}
		}
	}
	h.mu.RUnlock()
	writeJSON(w, http.StatusOK, map[string]bool{"allowed": allowed})
}
