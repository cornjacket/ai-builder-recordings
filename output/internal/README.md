# internal

Purpose: Top-level internal package index for the user service — contains the userservice sub-tree, which implements the full HTTP CRUD stack for user records.
The internal/ boundary ensures these packages are not importable outside the module.

Tags: architecture, overview

## Overview

The `internal/` directory contains a single composite sub-package, `userservice/`, which provides the complete user management implementation: HTTP routing, JSON encoding, and thread-safe in-memory storage. All application logic lives inside this boundary; `main.go` at the project root is the only caller that wires and starts the server.

## Components

| Sub-package | Responsibility | Doc |
|-------------|----------------|-----|
| [`userservice/`](userservice/README.md) | Composite package pairing an HTTP handler layer with a thread-safe in-memory store to form the full user CRUD service |

## Synthesis Docs

- [data-flow.md](data-flow.md) — how a request enters the internal boundary and flows through the userservice sub-tree
