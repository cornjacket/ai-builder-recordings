# userservice

Purpose: Provides an in-memory HTTP user management service with CRUD operations over a REST API.

Tags: architecture, design

A self-contained package implementing create, read, update, and delete operations for user records. It exposes an HTTP server on port 8080 with JSON request and response bodies.

## Components

| Component | Description |
|-----------|-------------|
| [store](store/README.md) | Thread-safe in-memory store for User records |
| [handlers](handlers/README.md) | HTTP CRUD handlers wired to the store |

See [theory-of-operation.md](theory-of-operation.md) for data-flow details.
