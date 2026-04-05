# integrate-user-service

## Goal

Wire all components into a cohesive unit and verify this level's acceptance criteria. Write main.go at the output root: instantiate the store, pass it to the handlers, register routes on net/http.ServeMux (POST /users, GET /users/{id}, PUT /users/{id}, DELETE /users/{id}), and call http.ListenAndServe(":8080", mux). End-to-end acceptance criteria: POST /users returns 201 with generated id; GET /users/{id} returns 200 with correct record or 404; PUT /users/{id} returns 200 with updated record or 404; DELETE /users/{id} returns 204 or 404; server listens on port 8080; all responses are JSON.

## Context

### Level 1 — 5594e1-0000-build-1


