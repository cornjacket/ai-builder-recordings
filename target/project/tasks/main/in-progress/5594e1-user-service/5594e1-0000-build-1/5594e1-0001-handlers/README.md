<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: handlers

## Goal

HTTP CRUD handlers for the user management API. Accepts a store interface and returns an http.Handler (or registers on a mux). Full API contract — POST /users {"name":string,"email":string} → 201 {"id":string,"name":string,"email":string}; GET /users/{id} → 200 {"id":string,"name":string,"email":string} or 404 {}; PUT /users/{id} {"name":string,"email":string} → 200 {"id":string,"name":string,"email":string} or 404 {}; DELETE /users/{id} → 204 (no body) or 404 {}. All request and response bodies are JSON. Returns 400 on malformed JSON request bodies.

## Context

### Level 1 — 5594e1-0000-build-1


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
