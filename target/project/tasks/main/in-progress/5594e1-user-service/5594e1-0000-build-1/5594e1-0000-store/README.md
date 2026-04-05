<!-- This file is managed by the ai-builder pipeline. Do not hand-edit. -->
# Task: store

## Goal

Thread-safe in-memory store for User records. User model: {"id": string (UUID, assigned on create), "name": string, "email": string}. Exposes: Create(name, email string) User; Get(id string) (User, bool); Update(id, name, email string) (User, bool); Delete(id string) bool. Uses sync.RWMutex internally. No external dependencies beyond stdlib.

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
