# api

Purpose: Full HTTP API contract for the user management handlers. Documents every endpoint, request/response schemas, and all status codes.

Tags: architecture, design

## Endpoints

### POST /users

Create a new user.

**Request body:**
```json
{"name": "string", "email": "string"}
```

**Responses:**
- `201 Created` — user created
  ```json
  {"id": "string", "name": "string", "email": "string"}
  ```
- `400 Bad Request` — malformed JSON request body (no body)

---

### GET /users/{id}

Retrieve a user by ID.

**Path parameter:** `id` — UUID string

**Responses:**
- `200 OK` — user found
  ```json
  {"id": "string", "name": "string", "email": "string"}
  ```
- `404 Not Found` — no user with that ID
  ```json
  {}
  ```

---

### PUT /users/{id}

Replace the name and email of an existing user.

**Path parameter:** `id` — UUID string

**Request body:**
```json
{"name": "string", "email": "string"}
```

**Responses:**
- `200 OK` — user updated
  ```json
  {"id": "string", "name": "string", "email": "string"}
  ```
- `404 Not Found` — no user with that ID
  ```json
  {}
  ```
- `400 Bad Request` — malformed JSON request body (no body)

---

### DELETE /users/{id}

Delete a user by ID.

**Path parameter:** `id` — UUID string

**Responses:**
- `204 No Content` — user deleted (no body)
- `404 Not Found` — no user with that ID
  ```json
  {}
  ```
