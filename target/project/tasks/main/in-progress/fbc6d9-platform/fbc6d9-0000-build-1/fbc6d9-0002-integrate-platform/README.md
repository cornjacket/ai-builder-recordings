<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: integrate-platform

## Goal

Wire the metrics and iam components into a single binary at cmd/platform/main.go. The binary starts the metrics HTTP listener on port 8081 and the IAM HTTP listener on port 8082 in separate goroutines, then blocks until a shutdown signal. There must be exactly one main package in the entire codebase.

End-to-end acceptance tests (net/http against live listeners started in TestMain or a helper) must cover all 12 endpoints verbatim from the spec:

Metrics listener (port 8081):
POST /events body:{"type":"click-mouse","userId":"<string>","payload":{}} → 201 {"id":"<string>","type":"<string>","userId":"<string>","payload":{}};
POST /events body:{"type":"submit-form","userId":"<string>","payload":{}} → 201 (same shape);
GET /events → 200 JSON array containing previously posted events.

IAM listener (port 8082):
POST /users body:{"username":"<string>","password":"<string>"} → 201 {"id":"<string>","username":"<string>"} (no password field);
GET /users/{id} → 200 {"id":"<string>","username":"<string>"} or 404;
DELETE /users/{id} → 200/204 on existing user or 404 on missing;
POST /auth/login body:{"username":"<string>","password":"<string>"} → 200 {"token":"<string>"};
POST /auth/logout header:Authorization:Bearer <token> → 200/204;
POST /roles body:{"name":"<string>","permissions":["<string>"]} → 201 {"id":"<string>","name":"<string>","permissions":["<string>"]};
GET /roles → 200 JSON array;
POST /users/{id}/roles body:{"roleId":"<string>"} → 200/201;
GET /users/{id}/roles → 200 JSON array;
POST /authz/check body:{"userId":"<string>","permission":"<string>"} → 200 {"allowed":<bool>}.

## Context

### Level 1 — fbc6d9-0000-build-1


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
