# authz API

Purpose: Defines the HTTP request and response contract for all five authz endpoints. Covers role creation, role listing, role assignment, user-role retrieval, and permission checking.

Tags: architecture, design

## POST /roles

Create a new role.

**Request body:**
```json
{"name": "string", "permissions": ["string"]}
```

**Responses:**
- `201 Created` — `{"id": "string", "name": "string", "permissions": ["string"]}`
- `400 Bad Request` — missing or empty `name` field, or body cannot be decoded

`permissions` may be an empty array. The response echoes the stored role including the generated `id`.

## GET /roles

List all roles.

**Responses:**
- `200 OK` — `[{"id": "string", "name": "string", "permissions": ["string"]}]`

Returns an empty JSON array `[]` when no roles exist. Never returns 404.

## POST /users/{id}/roles

Assign an existing role to a user.

**Path parameter:** `{id}` — the user ID (not validated against lifecycle; any string is accepted as a user identifier)

**Request body:**
```json
{"roleId": "string"}
```

**Responses:**
- `201 Created` — empty body; the role has been added to the user's role list
- `404 Not Found` — `roleId` does not refer to an existing role in RoleStore
- `400 Bad Request` — body cannot be decoded or `roleId` is empty

Assigning the same role to a user twice is idempotent — a second assignment still returns 201 but does not create a duplicate entry.

## GET /users/{id}/roles

Retrieve all roles assigned to a user.

**Path parameter:** `{id}` — the user ID

**Responses:**
- `200 OK` — `[{"id": "string", "name": "string", "permissions": ["string"]}]`

Returns an empty JSON array `[]` when the user has no assigned roles or the user ID is unknown. Never returns 404.

## POST /authz/check

Check whether a user has a specific permission.

**Request body:**
```json
{"userId": "string", "permission": "string"}
```

**Responses:**
- `200 OK` — `{"allowed": true}` or `{"allowed": false}`
- `400 Bad Request` — body cannot be decoded

Always returns 200. Returns `{"allowed": false}` when the user has no assigned roles, when none of the assigned roles exist in RoleStore, or when no role's permissions slice contains the requested permission string.
