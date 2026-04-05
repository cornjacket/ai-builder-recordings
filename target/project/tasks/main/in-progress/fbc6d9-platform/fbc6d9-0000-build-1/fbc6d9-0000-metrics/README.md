<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: metrics

## Goal

In-memory metrics service for recording and listing frontend user interaction events. Serves on port 8081.

HTTP API (verbatim from spec):
POST /events — body: {"type":"click-mouse"|"submit-form","userId":"<string>","payload":{}} → 201 with event object {"id":"<string>","type":"<string>","userId":"<string>","payload":{}};
GET /events → 200 JSON array of event objects, each with fields: id, type, userId, payload.

Data model: Event{id string, type string, userId string, payload map[string]interface{}}.
Store: thread-safe in-memory slice; IDs are generated UUIDs on POST.
Package layout: store.go (EventStore with Add/List), handlers.go (PostEvent, GetEvents), routes.go (returns http.Handler).
Unit tests must cover store and handler logic; handler tests use net/http/httptest.

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
