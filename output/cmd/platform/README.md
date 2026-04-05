# platform

Purpose: Binary entry point that wires together the metrics and IAM internal packages and serves each on a dedicated port.

Tags: main

## Overview

The `platform` binary is the single runnable entry point for the monolith. It starts the metrics HTTP service on `:8081` and the IAM HTTP service on `:8082`, running them concurrently in the same process.

## File Index

| File | Description | Doc |
|------|-------------|-----|
| `main.go` | Instantiates handlers and starts both HTTP servers | [main.go.md](main.go.md) |

## Notes

The metrics server is launched in a background goroutine; the IAM server blocks the main goroutine. A fatal error from either server (e.g. port already in use) will terminate the entire process via `log.Fatalf`. There is no graceful shutdown mechanism.
