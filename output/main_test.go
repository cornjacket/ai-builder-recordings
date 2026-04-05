// Package main tests the integrated user-service via an in-process HTTP test server.
//
// Purpose: End-to-end integration tests covering the full CRUD lifecycle against a real handler and store.
// Tags: implementation, main
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/handlers"
	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

func TestUserServiceCRUD(t *testing.T) {
	srv := httptest.NewServer(handlers.New(store.New()))
	defer srv.Close()

	base := srv.URL

	// 1. POST /users → 201, body has id, name, email
	body := `{"name":"Alice","email":"alice@example.com"}`
	resp, err := http.Post(base+"/users", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /users: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /users: expected 201, got %d", resp.StatusCode)
	}
	var created struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("POST /users: decode body: %v", err)
	}
	resp.Body.Close()
	if created.ID == "" {
		t.Fatal("POST /users: id is empty")
	}
	if created.Name != "Alice" {
		t.Fatalf("POST /users: name = %q, want Alice", created.Name)
	}
	if created.Email != "alice@example.com" {
		t.Fatalf("POST /users: email = %q, want alice@example.com", created.Email)
	}

	id := created.ID

	// 2. GET /users/{id} (known) → 200, correct fields
	resp, err = http.Get(base + "/users/" + id)
	if err != nil {
		t.Fatalf("GET /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/%s: expected 200, got %d", id, resp.StatusCode)
	}
	var got struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("GET /users/%s: decode body: %v", id, err)
	}
	resp.Body.Close()
	if got.ID != id || got.Name != "Alice" || got.Email != "alice@example.com" {
		t.Fatalf("GET /users/%s: unexpected body %+v", id, got)
	}

	// 3. PUT /users/{id} (known) → 200, updated fields
	putBody := `{"name":"Bob","email":"bob@example.com"}`
	req, _ := http.NewRequest(http.MethodPut, base+"/users/"+id, strings.NewReader(putBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT /users/%s: expected 200, got %d", id, resp.StatusCode)
	}
	var updated struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("PUT /users/%s: decode body: %v", id, err)
	}
	resp.Body.Close()
	if updated.Name != "Bob" || updated.Email != "bob@example.com" {
		t.Fatalf("PUT /users/%s: unexpected body %+v", id, updated)
	}

	// 4. GET /users/{id} (verify update stuck) → 200
	resp, err = http.Get(base + "/users/" + id)
	if err != nil {
		t.Fatalf("GET after PUT /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET after PUT /users/%s: expected 200, got %d", id, resp.StatusCode)
	}
	var afterUpdate struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	json.NewDecoder(resp.Body).Decode(&afterUpdate)
	resp.Body.Close()
	if afterUpdate.Name != "Bob" || afterUpdate.Email != "bob@example.com" {
		t.Fatalf("GET after PUT /users/%s: update did not persist: %+v", id, afterUpdate)
	}

	// 5. DELETE /users/{id} → 204
	req, _ = http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE /users/%s: expected 204, got %d", id, resp.StatusCode)
	}
	resp.Body.Close()

	// 6. GET /users/{id} (after delete) → 404
	resp, err = http.Get(base + "/users/" + id)
	if err != nil {
		t.Fatalf("GET after DELETE /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET after DELETE /users/%s: expected 404, got %d", id, resp.StatusCode)
	}
	resp.Body.Close()

	// 7. DELETE /users/{id} (after delete) → 404
	req, _ = http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE again /users/%s: %v", id, err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("DELETE again /users/%s: expected 404, got %d", id, resp.StatusCode)
	}
	resp.Body.Close()

	// 8. GET /users/nonexistent → 404
	resp, err = http.Get(base + "/users/nonexistent-id")
	if err != nil {
		t.Fatalf("GET nonexistent: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET nonexistent: expected 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// 9. PUT /users/nonexistent → 404
	req, _ = http.NewRequest(http.MethodPut, base+"/users/nonexistent-id", strings.NewReader(`{"name":"X","email":"x@x.com"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT nonexistent: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("PUT nonexistent: expected 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
