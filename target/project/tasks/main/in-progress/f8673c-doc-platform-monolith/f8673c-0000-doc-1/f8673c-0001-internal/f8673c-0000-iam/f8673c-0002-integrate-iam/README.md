# integrate-iam

## Goal

Cross-component synthesis: iam.go.md documenting the New() wiring function, and README.md (update if stale, leave if complete) for this directory

## Context

### Level 1 — f8673c-0000-doc-1

### Level 2 — f8673c-0001-internal
Internal service packages subtree — contains `metrics` (event store + HTTP handlers) and `iam` (lifecycle + authz) as composite sub-packages

### Level 3 — f8673c-0000-iam
IAM package — contains lifecycle and authz sub-packages alongside a parent-level iam.go; pre-existing README.md is present

