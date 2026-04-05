# handlers

## Goal

HTTP handlers for POST /events and GET /events — accepts a store.Storer interface and encodes/decodes JSON

## Context

### Level 1 — f8673c-0000-doc-1

### Level 2 — f8673c-0001-internal
Internal service packages subtree — contains `metrics` (event store + HTTP handlers) and `iam` (lifecycle + authz) as composite sub-packages

### Level 3 — f8673c-0001-metrics
Metrics package — contains handlers and store sub-packages alongside a parent-level metrics.go; pre-existing README.md is present

