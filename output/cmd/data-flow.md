# data-flow.md

Purpose: Describes how the platform binary wires together the internal service packages and routes incoming HTTP traffic to each.
At startup, `main()` obtains a handler from each internal package and binds it to a dedicated port; all subsequent data flow is request-scoped.

Tags: data-flow, architecture

## Startup Wiring

When the binary starts, `main()` calls into the two internal packages to obtain their `http.Handler` implementations. This is the only coupling point between `cmd` and `internal` — after construction, each server operates independently.

```
cmd/platform/main.go
        |
        |-- metrics.New() --> internal/metrics  (handler construction)
        |-- iam.New()     --> internal/iam       (handler construction)
        |
        +--> goroutine: http.ListenAndServe(":8081", metricsHandler)
        +--> (main goroutine): http.ListenAndServe(":8082", iamHandler)
```

## Request Routing

Once both servers are running, inbound HTTP requests are dispatched entirely within the respective internal package — `cmd` plays no further role in request handling.

```
Client
  |
  |-- :8081 --> metricsHandler (owned by internal/metrics)
  |-- :8082 --> iamHandler     (owned by internal/iam)
```

## Failure Propagation

A bind failure on either port causes `log.Fatalf` to terminate the whole process. Because the metrics server runs in a background goroutine, a runtime error from it will also kill the process via `log.Fatalf`. There is no partial-failure or graceful-shutdown path.
