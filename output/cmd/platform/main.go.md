# main.go

Purpose: Entry point for the platform binary — instantiates the metrics and IAM HTTP handlers and starts each on its own port concurrently.
The metrics service listens on :8081 and the IAM service on :8082. A fatal log terminates the process if either server fails to bind.

Tags: main

## Public API

This file defines only `main`; there is no exported API.

## Key Internals

### `main()`

```go
func main()
```

Calls `metrics.New()` and `iam.New()` to obtain their respective `http.Handler` implementations, then launches the metrics server in a goroutine so both servers run concurrently. The IAM server runs on the main goroutine; a fatal error from either server terminates the process.

## Dependencies

| Package | Role |
|---------|------|
| `github.com/cornjacket/platform/internal/iam` | Provides the IAM HTTP handler via `iam.New()` |
| `github.com/cornjacket/platform/internal/metrics` | Provides the metrics HTTP handler via `metrics.New()` |
| `log` (stdlib) | Startup logging and fatal error reporting |
| `net/http` (stdlib) | `http.ListenAndServe` for both servers |
