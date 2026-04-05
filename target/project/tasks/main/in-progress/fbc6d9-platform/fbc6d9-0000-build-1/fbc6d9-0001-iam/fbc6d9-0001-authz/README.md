# authz

## Goal

In-memory role store and RBAC permission-check handler for the IAM listener.

Endpoints:
POST /roles {"name":string,"permissions":[]string} → 201 {"id":string,"name":string,"permissions":[]string} or 400;
GET /roles → 200 [{"id":string,"name":string,"permissions":[]string}];
POST /users/{id}/roles {"roleId":string} → 201 or 404;
GET /users/{id}/roles → 200 [{"id":string,"name":string,"permissions":[]string}];
POST /authz/check {"userId":string,"permission":string} → 200 {"allowed":bool}.

In-memory stores: RoleStore map[id]→{id string, name string, permissions []string}, UserRoles map[userID]→[]roleID. Permission check walks the user's assigned roles and returns allowed:true if any role's permissions slice contains the requested permission. Exposes an http.Handler (or *http.ServeMux) for registration by integrate-iam.

## Context

### Level 1 — fbc6d9-0000-build-1

### Level 2 — fbc6d9-0001-iam
Identity and access management listener (port 8082) composed of two sub-components — lifecycle (user CRUD and authentication) and authz-rbac (roles and permission checks) — that together handle all ten IAM HTTP endpoints.

