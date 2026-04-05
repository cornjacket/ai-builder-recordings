<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: integrate-iam

## Goal

Wire lifecycle and authz handlers into a single http.ServeMux and expose it for the platform main.go. Create iam.go (package iam) in the output directory. It must instantiate lifecycle.New() and authz.New(), register their routes on a shared mux (route prefixes: /users and /auth go to lifecycle; /roles, /authz, and /users/ sub-paths for roles go to authz — use pattern matching carefully so /users/{id}/roles routes to authz while /users/{id} routes to lifecycle), and return the mux via a NewMux() function. No port binding here — the platform main.go calls http.ListenAndServe(":8082", iam.NewMux()). Task Level: INTERNAL.

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
