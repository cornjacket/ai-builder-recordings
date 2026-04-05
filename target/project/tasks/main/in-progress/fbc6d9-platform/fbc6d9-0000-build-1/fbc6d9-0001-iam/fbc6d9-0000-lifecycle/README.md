<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: lifecycle

## Goal

In-memory user store and token-based authentication handler for the IAM listener.

Endpoints:
POST /users {"username":string,"password":string} → 201 {"id":string,"username":string} or 400;
GET /users/{id} → 200 {"id":string,"username":string} or 404;
DELETE /users/{id} → 200 or 404;
POST /auth/login {"username":string,"password":string} → 200 {"token":string} or 401;
POST /auth/logout header Authorization:Bearer <token> → 200 or 401.

In-memory stores: UserStore map[id]→{id string, username string, passwordHash string}, TokenStore map[token]→userID string. IDs generated with crypto/rand or math/rand UUID. Passwords hashed with bcrypt or sha256. Exposes an http.Handler (or *http.ServeMux) for registration by integrate-iam.

## Context

### Level 1 — fbc6d9-0000-build-1

### Level 2 — fbc6d9-0001-iam
Identity and access management listener (port 8082) composed of two sub-components — lifecycle (user CRUD and authentication) and authz-rbac (roles and permission checks) — that together handle all ten IAM HTTP endpoints.

## Components

_To be completed by the ARCHITECT._

## Design

_To be completed by the ARCHITECT._

## Acceptance Criteria

_To be completed by the ARCHITECT._

## Test Command

_To be completed by the ARCHITECT._

## Suggested Tools

_To be completed by the ARCHITECT._

## Notes

_None._
