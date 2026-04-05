// Package main_test provides end-to-end acceptance tests for the platform binary.
//
// Purpose: Spins up in-process httptest servers for metrics and IAM listeners and runs net/http
// tests against all 12 spec endpoints. Uses TestMain for shared server lifecycle.
// Tags: implementation, platform
package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	iampkg "github.com/cornjacket/platform-monolith/internal/iam"
	"github.com/cornjacket/platform-monolith/internal/metrics"
)

var (
	metricsBase string
	iamBase     string
)

func TestMain(m *testing.M) {
	mSrv := httptest.NewServer(metrics.NewRouter(metrics.NewEventStore()))
	iSrv := httptest.NewServer(iampkg.NewMux())
	metricsBase = mSrv.URL
	iamBase = iSrv.URL
	code := m.Run()
	mSrv.Close()
	iSrv.Close()
	os.Exit(code)
}

// mustDo executes an HTTP request and fails the test on transport error.
func mustDo(t *testing.T, req *http.Request) *http.Response {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

// jsonBody encodes v to a *bytes.Reader suitable for http.NewRequest.
func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

// TestPostEventClickMouse verifies POST /events with type=click-mouse returns 201 with correct fields.
func TestPostEventClickMouse(t *testing.T) {
	body := jsonBody(map[string]any{"type": "click-mouse", "userId": "u1", "payload": map[string]any{}})
	req, _ := http.NewRequest(http.MethodPost, metricsBase+"/events", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got["id"] == "" || got["id"] == nil {
		t.Error("expected non-empty id")
	}
	if got["type"] != "click-mouse" {
		t.Errorf("expected type=click-mouse, got %v", got["type"])
	}
	if got["userId"] != "u1" {
		t.Errorf("expected userId=u1, got %v", got["userId"])
	}
	if _, ok := got["payload"]; !ok {
		t.Error("expected payload field")
	}
}

// TestPostEventSubmitForm verifies POST /events with type=submit-form returns 201.
func TestPostEventSubmitForm(t *testing.T) {
	body := jsonBody(map[string]any{"type": "submit-form", "userId": "u2", "payload": map[string]any{}})
	req, _ := http.NewRequest(http.MethodPost, metricsBase+"/events", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got["type"] != "submit-form" {
		t.Errorf("expected type=submit-form, got %v", got["type"])
	}
}

// TestGetEvents verifies GET /events returns 200 and a JSON array containing posted events.
func TestGetEvents(t *testing.T) {
	// Post two events first so we have known IDs.
	var ids []string
	for _, evType := range []string{"click-mouse", "submit-form"} {
		body := jsonBody(map[string]any{"type": evType, "userId": "u-ge", "payload": map[string]any{}})
		req, _ := http.NewRequest(http.MethodPost, metricsBase+"/events", body)
		req.Header.Set("Content-Type", "application/json")
		resp := mustDo(t, req)
		var got map[string]any
		json.NewDecoder(resp.Body).Decode(&got) //nolint:errcheck
		resp.Body.Close()
		if id, ok := got["id"].(string); ok && id != "" {
			ids = append(ids, id)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, metricsBase+"/events", nil)
	resp := mustDo(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var events []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	idSet := make(map[string]bool)
	for _, e := range events {
		if id, ok := e["id"].(string); ok {
			idSet[id] = true
		}
	}
	for _, id := range ids {
		if !idSet[id] {
			t.Errorf("event id %s not found in GET /events response", id)
		}
	}
}

// TestPostUser verifies POST /users returns 201 with id and username but no password field.
func TestPostUser(t *testing.T) {
	body := jsonBody(map[string]string{"username": "alice", "password": "secret"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if id, ok := got["id"].(string); !ok || id == "" {
		t.Error("expected non-empty id")
	}
	if got["username"] != "alice" {
		t.Errorf("expected username=alice, got %v", got["username"])
	}
	if _, hasPass := got["password"]; hasPass {
		t.Error("response must not contain password field")
	}
}

// TestGetUser verifies GET /users/{id} returns 200 for existing and 404 for missing.
func TestGetUser(t *testing.T) {
	// Create a user.
	body := jsonBody(map[string]string{"username": "bob", "password": "pw"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created) //nolint:errcheck
	resp.Body.Close()

	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("could not get created user id")
	}

	// GET existing.
	req2, _ := http.NewRequest(http.MethodGet, iamBase+"/users/"+id, nil)
	resp2 := mustDo(t, req2)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for existing user, got %d", resp2.StatusCode)
	}
	var got map[string]any
	json.NewDecoder(resp2.Body).Decode(&got) //nolint:errcheck
	if got["id"] != id {
		t.Errorf("expected id=%s, got %v", id, got["id"])
	}

	// GET non-existent.
	req3, _ := http.NewRequest(http.MethodGet, iamBase+"/users/nonexistent-id", nil)
	resp3 := mustDo(t, req3)
	resp3.Body.Close()
	if resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing user, got %d", resp3.StatusCode)
	}
}

// TestDeleteUser verifies DELETE /users/{id} returns 200/204 for existing and 404 for missing.
func TestDeleteUser(t *testing.T) {
	// Create a user to delete.
	body := jsonBody(map[string]string{"username": "charlie", "password": "pw"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created) //nolint:errcheck
	resp.Body.Close()

	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("could not get created user id")
	}

	// DELETE existing.
	req2, _ := http.NewRequest(http.MethodDelete, iamBase+"/users/"+id, nil)
	resp2 := mustDo(t, req2)
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 200 or 204 for existing user delete, got %d", resp2.StatusCode)
	}

	// DELETE non-existent.
	req3, _ := http.NewRequest(http.MethodDelete, iamBase+"/users/nonexistent-id", nil)
	resp3 := mustDo(t, req3)
	resp3.Body.Close()
	if resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing user delete, got %d", resp3.StatusCode)
	}
}

// TestLogin verifies POST /auth/login returns 200 with a non-empty token.
func TestLogin(t *testing.T) {
	// Create user.
	body := jsonBody(map[string]string{"username": "dave", "password": "pass"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	resp.Body.Close()

	// Login.
	body2 := jsonBody(map[string]string{"username": "dave", "password": "pass"})
	req2, _ := http.NewRequest(http.MethodPost, iamBase+"/auth/login", body2)
	req2.Header.Set("Content-Type", "application/json")
	resp2 := mustDo(t, req2)
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	token, ok := got["token"].(string)
	if !ok || token == "" {
		t.Error("expected non-empty token")
	}
}

// TestLogout verifies POST /auth/logout with a valid bearer token returns 200/204.
func TestLogout(t *testing.T) {
	// Create user and login.
	body := jsonBody(map[string]string{"username": "eve", "password": "pass"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	mustDo(t, req).Body.Close()

	body2 := jsonBody(map[string]string{"username": "eve", "password": "pass"})
	req2, _ := http.NewRequest(http.MethodPost, iamBase+"/auth/login", body2)
	req2.Header.Set("Content-Type", "application/json")
	resp2 := mustDo(t, req2)
	var loginResp map[string]any
	json.NewDecoder(resp2.Body).Decode(&loginResp) //nolint:errcheck
	resp2.Body.Close()

	token, _ := loginResp["token"].(string)
	if token == "" {
		t.Fatal("could not get token")
	}

	// Logout.
	req3, _ := http.NewRequest(http.MethodPost, iamBase+"/auth/logout", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	resp3 := mustDo(t, req3)
	resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK && resp3.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 200 or 204, got %d", resp3.StatusCode)
	}
}

// TestPostRole verifies POST /roles returns 201 with id, name, and permissions fields.
func TestPostRole(t *testing.T) {
	body := jsonBody(map[string]any{"name": "admin", "permissions": []string{"read", "write"}})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/roles", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if id, ok := got["id"].(string); !ok || id == "" {
		t.Error("expected non-empty id")
	}
	if got["name"] != "admin" {
		t.Errorf("expected name=admin, got %v", got["name"])
	}
	perms, ok := got["permissions"].([]any)
	if !ok || len(perms) < 2 {
		t.Errorf("expected permissions array with at least 2 entries, got %v", got["permissions"])
	}
}

// TestGetRoles verifies GET /roles returns 200 and a JSON array containing created roles.
func TestGetRoles(t *testing.T) {
	// Create a role.
	body := jsonBody(map[string]any{"name": "viewer", "permissions": []string{"read"}})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/roles", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created) //nolint:errcheck
	resp.Body.Close()

	roleID, _ := created["id"].(string)

	// List roles.
	req2, _ := http.NewRequest(http.MethodGet, iamBase+"/roles", nil)
	resp2 := mustDo(t, req2)
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
	var roles []map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&roles); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	found := false
	for _, r := range roles {
		if r["id"] == roleID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created role %s not found in GET /roles", roleID)
	}
}

// TestAssignRole verifies POST /users/{id}/roles returns 200/201.
func TestAssignRole(t *testing.T) {
	// Create user.
	body := jsonBody(map[string]string{"username": "frank", "password": "pw"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var user map[string]any
	json.NewDecoder(resp.Body).Decode(&user) //nolint:errcheck
	resp.Body.Close()
	userID, _ := user["id"].(string)

	// Create role.
	body2 := jsonBody(map[string]any{"name": "editor", "permissions": []string{"write"}})
	req2, _ := http.NewRequest(http.MethodPost, iamBase+"/roles", body2)
	req2.Header.Set("Content-Type", "application/json")
	resp2 := mustDo(t, req2)
	var role map[string]any
	json.NewDecoder(resp2.Body).Decode(&role) //nolint:errcheck
	resp2.Body.Close()
	roleID, _ := role["id"].(string)

	if userID == "" || roleID == "" {
		t.Fatal("could not get user or role id")
	}

	// Assign role.
	body3 := jsonBody(map[string]string{"roleId": roleID})
	req3, _ := http.NewRequest(http.MethodPost, iamBase+"/users/"+userID+"/roles", body3)
	req3.Header.Set("Content-Type", "application/json")
	resp3 := mustDo(t, req3)
	resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK && resp3.StatusCode != http.StatusCreated {
		t.Fatalf("expected 200 or 201, got %d", resp3.StatusCode)
	}
}

// TestGetUserRoles verifies GET /users/{id}/roles returns 200 and a JSON array with assigned roles.
func TestGetUserRoles(t *testing.T) {
	// Create user.
	body := jsonBody(map[string]string{"username": "grace", "password": "pw"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var user map[string]any
	json.NewDecoder(resp.Body).Decode(&user) //nolint:errcheck
	resp.Body.Close()
	userID, _ := user["id"].(string)

	// Create role.
	body2 := jsonBody(map[string]any{"name": "moderator", "permissions": []string{"read", "delete"}})
	req2, _ := http.NewRequest(http.MethodPost, iamBase+"/roles", body2)
	req2.Header.Set("Content-Type", "application/json")
	resp2 := mustDo(t, req2)
	var role map[string]any
	json.NewDecoder(resp2.Body).Decode(&role) //nolint:errcheck
	resp2.Body.Close()
	roleID, _ := role["id"].(string)

	if userID == "" || roleID == "" {
		t.Fatal("could not get user or role id")
	}

	// Assign role.
	body3 := jsonBody(map[string]string{"roleId": roleID})
	req3, _ := http.NewRequest(http.MethodPost, iamBase+"/users/"+userID+"/roles", body3)
	req3.Header.Set("Content-Type", "application/json")
	mustDo(t, req3).Body.Close()

	// Get user roles.
	req4, _ := http.NewRequest(http.MethodGet, iamBase+"/users/"+userID+"/roles", nil)
	resp4 := mustDo(t, req4)
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp4.StatusCode)
	}
	var roles []map[string]any
	if err := json.NewDecoder(resp4.Body).Decode(&roles); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	found := false
	for _, r := range roles {
		if r["id"] == roleID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("assigned role %s not found in GET /users/%s/roles", roleID, userID)
	}
}

// TestAuthzCheck verifies POST /authz/check returns allowed=true for held permissions and allowed=false otherwise.
func TestAuthzCheck(t *testing.T) {
	// Create user.
	body := jsonBody(map[string]string{"username": "henry", "password": "pw"})
	req, _ := http.NewRequest(http.MethodPost, iamBase+"/users", body)
	req.Header.Set("Content-Type", "application/json")
	resp := mustDo(t, req)
	var user map[string]any
	json.NewDecoder(resp.Body).Decode(&user) //nolint:errcheck
	resp.Body.Close()
	userID, _ := user["id"].(string)

	// Create role with "read" permission.
	body2 := jsonBody(map[string]any{"name": "reader", "permissions": []string{"read"}})
	req2, _ := http.NewRequest(http.MethodPost, iamBase+"/roles", body2)
	req2.Header.Set("Content-Type", "application/json")
	resp2 := mustDo(t, req2)
	var role map[string]any
	json.NewDecoder(resp2.Body).Decode(&role) //nolint:errcheck
	resp2.Body.Close()
	roleID, _ := role["id"].(string)

	if userID == "" || roleID == "" {
		t.Fatal("could not get user or role id")
	}

	// Assign role.
	body3 := jsonBody(map[string]string{"roleId": roleID})
	req3, _ := http.NewRequest(http.MethodPost, iamBase+"/users/"+userID+"/roles", body3)
	req3.Header.Set("Content-Type", "application/json")
	mustDo(t, req3).Body.Close()

	// Check allowed permission.
	body4 := jsonBody(map[string]string{"userId": userID, "permission": "read"})
	req4, _ := http.NewRequest(http.MethodPost, iamBase+"/authz/check", body4)
	req4.Header.Set("Content-Type", "application/json")
	resp4 := mustDo(t, req4)
	defer resp4.Body.Close()

	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp4.StatusCode)
	}
	var result map[string]any
	json.NewDecoder(resp4.Body).Decode(&result) //nolint:errcheck
	if result["allowed"] != true {
		t.Errorf("expected allowed=true for held permission, got %v", result["allowed"])
	}

	// Check disallowed permission.
	body5 := jsonBody(map[string]string{"userId": userID, "permission": "delete"})
	req5, _ := http.NewRequest(http.MethodPost, iamBase+"/authz/check", body5)
	req5.Header.Set("Content-Type", "application/json")
	resp5 := mustDo(t, req5)
	defer resp5.Body.Close()

	var result2 map[string]any
	json.NewDecoder(resp5.Body).Decode(&result2) //nolint:errcheck
	if result2["allowed"] != false {
		t.Errorf("expected allowed=false for non-held permission, got %v", result2["allowed"])
	}
}
