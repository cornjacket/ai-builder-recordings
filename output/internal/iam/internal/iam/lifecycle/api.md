# lifecycle API

Purpose: Defines the HTTP request and response contract for all five lifecycle endpoints. Covers user registration, retrieval, deletion, login, and logout.

Tags: architecture, design

## POST /users

Register a new user.

**Request body:**
```json
{"username": "string", "password": "string"}
```

**Responses:**
- `201 Created` — `{"id": "string", "username": "string"}`
- `400 Bad Request` — missing fields, empty values, or username already taken

The response body never includes the password or password hash.

## GET /users/{id}

Retrieve a user by ID.

**Responses:**
- `200 OK` — `{"id": "string", "username": "string"}`
- `404 Not Found` — no user with that ID

## DELETE /users/{id}

Delete a user and invalidate all tokens belonging to that user.

**Responses:**
- `200 OK` — empty body
- `404 Not Found` — no user with that ID

## POST /auth/login

Authenticate a user and issue an opaque session token.

**Request body:**
```json
{"username": "string", "password": "string"}
```

**Responses:**
- `200 OK` — `{"token": "string"}`
- `401 Unauthorized` — unknown username or wrong password

## POST /auth/logout

Invalidate a session token.

**Request header:** `Authorization: Bearer <token>`

**Responses:**
- `200 OK` — empty body
- `401 Unauthorized` — missing `Authorization` header, wrong scheme, or unknown/already-invalidated token
