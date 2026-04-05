// Package lifecycle provides in-memory user CRUD and token-based authentication for the IAM listener.
//
// Purpose: Table-driven HTTP tests for all five lifecycle endpoints.
// Tags: implementation, lifecycle
package lifecycle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHandler() *http.ServeMux {
	return Handler()
}

func doRequest(t *testing.T, mux http.Handler, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body != "" {
		buf = bytes.NewBufferString(body)
	} else {
		buf = &bytes.Buffer{}
	}
	req := httptest.NewRequest(method, path, buf)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func TestPostUsers_Success(t *testing.T) {
	mux := newTestHandler()
	rr := doRequest(t, mux, http.MethodPost, "/users", `{"username":"alice","password":"secret"}`, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id")
	}
	if resp["username"] != "alice" {
		t.Errorf("expected username alice got %q", resp["username"])
	}
	if _, ok := resp["password"]; ok {
		t.Error("response must not contain password field")
	}
	if _, ok := resp["passwordHash"]; ok {
		t.Error("response must not contain passwordHash field")
	}
}

func TestPostUsers_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"missing password", `{"username":"bob"}`},
		{"missing username", `{"password":"secret"}`},
		{"empty username", `{"username":"","password":"secret"}`},
		{"empty password", `{"username":"bob","password":""}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := newTestHandler()
			rr := doRequest(t, mux, http.MethodPost, "/users", tc.body, nil)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected 400 got %d", rr.Code)
			}
		})
	}
}

func TestPostUsers_DuplicateUsername(t *testing.T) {
	mux := newTestHandler()
	doRequest(t, mux, http.MethodPost, "/users", `{"username":"alice","password":"secret"}`, nil)
	rr := doRequest(t, mux, http.MethodPost, "/users", `{"username":"alice","password":"other"}`, nil)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", rr.Code)
	}
}

func createUser(t *testing.T, mux http.Handler, username, password string) string {
	t.Helper()
	rr := doRequest(t, mux, http.MethodPost, "/users", `{"username":"`+username+`","password":"`+password+`"}`, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("createUser: expected 201 got %d body=%s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck
	return resp["id"]
}

func TestGetUser_Found(t *testing.T) {
	mux := newTestHandler()
	id := createUser(t, mux, "alice", "secret")

	rr := doRequest(t, mux, http.MethodGet, "/users/"+id, "", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck
	if resp["id"] != id {
		t.Errorf("expected id %q got %q", id, resp["id"])
	}
	if resp["username"] != "alice" {
		t.Errorf("expected username alice got %q", resp["username"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mux := newTestHandler()
	rr := doRequest(t, mux, http.MethodGet, "/users/nonexistent-id", "", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d", rr.Code)
	}
}

func TestDeleteUser_Existing(t *testing.T) {
	mux := newTestHandler()
	id := createUser(t, mux, "alice", "secret")

	rr := doRequest(t, mux, http.MethodDelete, "/users/"+id, "", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rr.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mux := newTestHandler()
	rr := doRequest(t, mux, http.MethodDelete, "/users/nonexistent-id", "", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d", rr.Code)
	}
}

func TestDeleteUser_ThenGet(t *testing.T) {
	mux := newTestHandler()
	id := createUser(t, mux, "alice", "secret")
	doRequest(t, mux, http.MethodDelete, "/users/"+id, "", nil)

	rr := doRequest(t, mux, http.MethodGet, "/users/"+id, "", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete got %d", rr.Code)
	}
}

func loginUser(t *testing.T, mux http.Handler, username, password string) (string, int) {
	t.Helper()
	body := `{"username":"` + username + `","password":"` + password + `"}`
	rr := doRequest(t, mux, http.MethodPost, "/auth/login", body, nil)
	if rr.Code != http.StatusOK {
		return "", rr.Code
	}
	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck
	return resp["token"], rr.Code
}

func TestLogin_Success(t *testing.T) {
	mux := newTestHandler()
	createUser(t, mux, "alice", "secret")

	token, code := loginUser(t, mux, "alice", "secret")
	if code != http.StatusOK {
		t.Fatalf("expected 200 got %d", code)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	mux := newTestHandler()
	createUser(t, mux, "alice", "secret")

	_, code := loginUser(t, mux, "alice", "wrongpassword")
	if code != http.StatusUnauthorized {
		t.Errorf("expected 401 got %d", code)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	mux := newTestHandler()
	_, code := loginUser(t, mux, "nobody", "secret")
	if code != http.StatusUnauthorized {
		t.Errorf("expected 401 got %d", code)
	}
}

func TestLogout_Success(t *testing.T) {
	mux := newTestHandler()
	createUser(t, mux, "alice", "secret")
	token, _ := loginUser(t, mux, "alice", "secret")

	rr := doRequest(t, mux, http.MethodPost, "/auth/logout", "", map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rr.Code)
	}
}

func TestLogout_Reuse(t *testing.T) {
	mux := newTestHandler()
	createUser(t, mux, "alice", "secret")
	token, _ := loginUser(t, mux, "alice", "secret")

	doRequest(t, mux, http.MethodPost, "/auth/logout", "", map[string]string{
		"Authorization": "Bearer " + token,
	})
	rr := doRequest(t, mux, http.MethodPost, "/auth/logout", "", map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 on second logout got %d", rr.Code)
	}
}

func TestLogout_NoHeader(t *testing.T) {
	mux := newTestHandler()
	rr := doRequest(t, mux, http.MethodPost, "/auth/logout", "", nil)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 got %d", rr.Code)
	}
}

func TestLogout_WrongScheme(t *testing.T) {
	mux := newTestHandler()
	createUser(t, mux, "alice", "secret")
	token, _ := loginUser(t, mux, "alice", "secret")

	rr := doRequest(t, mux, http.MethodPost, "/auth/logout", "", map[string]string{
		"Authorization": "Token " + token,
	})
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 got %d", rr.Code)
	}
}

func TestDeleteUser_CascadesToken(t *testing.T) {
	mux := newTestHandler()
	id := createUser(t, mux, "alice", "secret")
	token, _ := loginUser(t, mux, "alice", "secret")

	doRequest(t, mux, http.MethodDelete, "/users/"+id, "", nil)

	rr := doRequest(t, mux, http.MethodPost, "/auth/logout", "", map[string]string{
		"Authorization": "Bearer " + token,
	})
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 after user deletion got %d", rr.Code)
	}
}

func TestConcurrentRequests(t *testing.T) {
	mux := newTestHandler()
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func(n int) {
			defer func() { done <- struct{}{} }()
			username := "user" + string(rune('a'+n))
			id := createUser(t, mux, username, "pass")
			loginUser(t, mux, username, "pass")       //nolint:errcheck
			doRequest(t, mux, http.MethodGet, "/users/"+id, "", nil)
		}(i)
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
