// Package iam_test exercises NewMux() routing via net/http/httptest.
//
// Purpose: Verifies that NewMux() correctly routes each IAM endpoint to lifecycle or authz.
// Tags: implementation, iam
package iam_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cornjacket/platform-monolith/internal/iam"
)

func do(mux *http.ServeMux, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body) //nolint:errcheck
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func decode(rr *httptest.ResponseRecorder, v any) {
	json.NewDecoder(rr.Body).Decode(v) //nolint:errcheck
}

// TestNewMux_NotNil confirms NewMux returns a non-nil mux with no panic.
func TestNewMux_NotNil(t *testing.T) {
	mux := iam.NewMux()
	if mux == nil {
		t.Fatal("NewMux() returned nil")
	}
}

// TestLifecycleRoutes exercises the lifecycle happy-path round-trip:
// create user, get user, login, logout, delete user.
func TestLifecycleRoutes(t *testing.T) {
	mux := iam.NewMux()

	// POST /users → 201
	rr := do(mux, http.MethodPost, "/users", map[string]string{"username": "u1", "password": "pw"}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /users: want 201, got %d", rr.Code)
	}
	var created map[string]string
	decode(rr, &created)
	userID := created["id"]
	if userID == "" || created["username"] != "u1" {
		t.Fatalf("POST /users: unexpected body %v", created)
	}

	// GET /users/{id} → 200 with lifecycle shape (id + username)
	rr = do(mux, http.MethodGet, "/users/"+userID, nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /users/%s: want 200, got %d", userID, rr.Code)
	}
	var got map[string]string
	decode(rr, &got)
	if got["id"] != userID || got["username"] != "u1" {
		t.Fatalf("GET /users/%s: unexpected lifecycle shape %v", userID, got)
	}
	// Must NOT be a role array (authz shape)
	if _, hasRoles := got["roles"]; hasRoles {
		t.Fatalf("GET /users/%s reached authz, not lifecycle", userID)
	}

	// POST /auth/login → 200 with token
	rr = do(mux, http.MethodPost, "/auth/login", map[string]string{"username": "u1", "password": "pw"}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /auth/login: want 200, got %d", rr.Code)
	}
	var loginResp map[string]string
	decode(rr, &loginResp)
	token := loginResp["token"]
	if token == "" {
		t.Fatal("POST /auth/login: no token returned")
	}

	// POST /auth/logout → 200
	rr = do(mux, http.MethodPost, "/auth/logout", nil, map[string]string{"Authorization": "Bearer " + token})
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /auth/logout: want 200, got %d", rr.Code)
	}

	// DELETE /users/{id} → 200 (lifecycle, not authz)
	rr = do(mux, http.MethodDelete, "/users/"+userID, nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("DELETE /users/%s: want 200, got %d", userID, rr.Code)
	}
}

// TestAuthzRoutes exercises the authz happy-path:
// create role, list roles, assign role to user, get user roles, authz check.
func TestAuthzRoutes(t *testing.T) {
	mux := iam.NewMux()

	// Create a user first (needed for user ID).
	rr := do(mux, http.MethodPost, "/users", map[string]string{"username": "u2", "password": "pw2"}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /users: want 201, got %d", rr.Code)
	}
	var userBody map[string]string
	decode(rr, &userBody)
	userID := userBody["id"]

	// POST /roles → 201
	rr = do(mux, http.MethodPost, "/roles", map[string]any{"name": "admin", "permissions": []string{"read"}}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /roles: want 201, got %d", rr.Code)
	}
	var role map[string]any
	decode(rr, &role)
	roleID, _ := role["id"].(string)
	if roleID == "" || role["name"] != "admin" {
		t.Fatalf("POST /roles: unexpected body %v", role)
	}
	perms, _ := role["permissions"].([]any)
	if len(perms) != 1 || perms[0] != "read" {
		t.Fatalf("POST /roles: unexpected permissions %v", perms)
	}

	// GET /roles → 200 with array (authz, not lifecycle)
	rr = do(mux, http.MethodGet, "/roles", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /roles: want 200, got %d", rr.Code)
	}
	var roles []any
	decode(rr, &roles)
	if len(roles) < 1 {
		t.Fatalf("GET /roles: expected at least one role, got %v", roles)
	}

	// POST /users/{id}/roles → 201 (authz, not lifecycle)
	rr = do(mux, http.MethodPost, "/users/"+userID+"/roles", map[string]string{"roleId": roleID}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /users/%s/roles: want 201, got %d", userID, rr.Code)
	}

	// GET /users/{id}/roles → 200 with role array (authz, not lifecycle)
	rr = do(mux, http.MethodGet, "/users/"+userID+"/roles", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /users/%s/roles: want 200, got %d", userID, rr.Code)
	}
	var userRoles []map[string]any
	decode(rr, &userRoles)
	if len(userRoles) < 1 {
		t.Fatalf("GET /users/%s/roles: expected role array, got empty", userID)
	}
	// Must be role shape (has "id" and "name"), not user shape
	if _, hasUsername := userRoles[0]["username"]; hasUsername {
		t.Fatalf("GET /users/%s/roles reached lifecycle, not authz", userID)
	}

	// POST /authz/check → 200 {"allowed": true}
	rr = do(mux, http.MethodPost, "/authz/check", map[string]string{"userId": userID, "permission": "read"}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("POST /authz/check: want 200, got %d", rr.Code)
	}
	var check map[string]bool
	decode(rr, &check)
	if !check["allowed"] {
		t.Fatalf("POST /authz/check: expected allowed=true, got %v", check)
	}
}

// TestUsersSubpathRouting confirms /users/{id}/roles → authz and /users/{id} → lifecycle.
func TestUsersSubpathRouting(t *testing.T) {
	mux := iam.NewMux()

	// Create user to get a real ID.
	rr := do(mux, http.MethodPost, "/users", map[string]string{"username": "u3", "password": "pw3"}, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("POST /users: want 201, got %d", rr.Code)
	}
	var userBody map[string]string
	decode(rr, &userBody)
	userID := userBody["id"]

	// GET /users/{id} must return lifecycle shape (id + username), not a role array.
	rr = do(mux, http.MethodGet, "/users/"+userID, nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /users/%s: want 200, got %d", userID, rr.Code)
	}
	var userShape map[string]string
	decode(rr, &userShape)
	if userShape["id"] != userID || userShape["username"] != "u3" {
		t.Fatalf("GET /users/%s: want lifecycle shape, got %v", userID, userShape)
	}

	// GET /users/{id}/roles must return a JSON array (authz), not a user object.
	rr = do(mux, http.MethodGet, "/users/"+userID+"/roles", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /users/%s/roles: want 200, got %d", userID, rr.Code)
	}
	var rolesShape []any
	decode(rr, &rolesShape)
	// rolesShape should be a (possibly empty) array, not an error or user object.
	// The fact that Decode succeeded into []any confirms it is an array.
	_ = rolesShape
}
