// Package authz provides in-memory RBAC role management and permission-check handling.
//
// Purpose: Tests for all five authz endpoints using httptest.
// Tags: implementation, authz
package authz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newServer() *http.ServeMux {
	return Handler()
}

func doRequest(t *testing.T, mux *http.ServeMux, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func decodeJSON(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

// createRole is a test helper that creates a role and returns its ID.
func createRole(t *testing.T, mux *http.ServeMux, name string, perms []string) string {
	t.Helper()
	rr := doRequest(t, mux, http.MethodPost, "/roles", map[string]any{
		"name":        name,
		"permissions": perms,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("createRole: expected 201, got %d", rr.Code)
	}
	var role Role
	decodeJSON(t, rr, &role)
	if role.ID == "" {
		t.Fatal("createRole: expected non-empty id")
	}
	return role.ID
}

func TestCreateRole_Success(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodPost, "/roles", map[string]any{
		"name":        "admin",
		"permissions": []string{"read", "write"},
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var role Role
	decodeJSON(t, rr, &role)
	if role.Name != "admin" {
		t.Errorf("expected name=admin, got %q", role.Name)
	}
	if role.ID == "" {
		t.Error("expected non-empty id")
	}
	if len(role.Permissions) != 2 || role.Permissions[0] != "read" || role.Permissions[1] != "write" {
		t.Errorf("unexpected permissions: %v", role.Permissions)
	}
}

func TestCreateRole_MissingName(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodPost, "/roles", map[string]any{})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestListRoles_Empty(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodGet, "/roles", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var roles []Role
	decodeJSON(t, rr, &roles)
	if len(roles) != 0 {
		t.Errorf("expected empty array, got %d elements", len(roles))
	}
}

func TestListRoles_TwoRoles(t *testing.T) {
	mux := newServer()
	id1 := createRole(t, mux, "admin", []string{"read", "write"})
	id2 := createRole(t, mux, "viewer", []string{"read"})

	rr := doRequest(t, mux, http.MethodGet, "/roles", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var roles []Role
	decodeJSON(t, rr, &roles)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	ids := map[string]bool{roles[0].ID: true, roles[1].ID: true}
	if !ids[id1] || !ids[id2] {
		t.Errorf("unexpected role IDs: %v", ids)
	}
}

func TestAssignRole_Success(t *testing.T) {
	mux := newServer()
	roleID := createRole(t, mux, "admin", []string{"read"})
	rr := doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestAssignRole_NonexistentRole(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": "nonexistent-id"})
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestGetUserRoles_OneRole(t *testing.T) {
	mux := newServer()
	roleID := createRole(t, mux, "admin", []string{"read", "write"})
	doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})

	rr := doRequest(t, mux, http.MethodGet, "/users/u1/roles", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var roles []Role
	decodeJSON(t, rr, &roles)
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if roles[0].ID != roleID {
		t.Errorf("expected role id %q, got %q", roleID, roles[0].ID)
	}
}

func TestGetUserRoles_NoRoles(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodGet, "/users/u1/roles", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var roles []Role
	decodeJSON(t, rr, &roles)
	if len(roles) != 0 {
		t.Errorf("expected empty array, got %d elements", len(roles))
	}
}

func TestAuthzCheck_Allowed(t *testing.T) {
	mux := newServer()
	roleID := createRole(t, mux, "admin", []string{"read", "write"})
	doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})

	rr := doRequest(t, mux, http.MethodPost, "/authz/check", map[string]string{
		"userId":     "u1",
		"permission": "read",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result map[string]bool
	decodeJSON(t, rr, &result)
	if !result["allowed"] {
		t.Error("expected allowed=true")
	}
}

func TestAuthzCheck_Denied(t *testing.T) {
	mux := newServer()
	roleID := createRole(t, mux, "admin", []string{"read", "write"})
	doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})

	rr := doRequest(t, mux, http.MethodPost, "/authz/check", map[string]string{
		"userId":     "u1",
		"permission": "delete",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result map[string]bool
	decodeJSON(t, rr, &result)
	if result["allowed"] {
		t.Error("expected allowed=false")
	}
}

func TestAuthzCheck_UnknownUser(t *testing.T) {
	mux := newServer()
	rr := doRequest(t, mux, http.MethodPost, "/authz/check", map[string]string{
		"userId":     "unknown",
		"permission": "read",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result map[string]bool
	decodeJSON(t, rr, &result)
	if result["allowed"] {
		t.Error("expected allowed=false for unknown user")
	}
}

func TestAssignRole_Idempotent(t *testing.T) {
	mux := newServer()
	roleID := createRole(t, mux, "admin", []string{"read"})
	doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})
	doRequest(t, mux, http.MethodPost, "/users/u1/roles", map[string]string{"roleId": roleID})

	rr := doRequest(t, mux, http.MethodGet, "/users/u1/roles", nil)
	var roles []Role
	decodeJSON(t, rr, &roles)
	if len(roles) != 1 {
		t.Errorf("expected role to appear exactly once, got %d", len(roles))
	}
}
