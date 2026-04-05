// Package main is the platform binary entry point.
//
// Purpose: Starts the metrics HTTP listener on port 8081 and the IAM HTTP listener on port 8082 in
// separate goroutines, then blocks until SIGINT or SIGTERM triggers a graceful shutdown.
// Tags: implementation, platform
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	iampkg "github.com/cornjacket/platform-monolith/internal/iam"
	"github.com/cornjacket/platform-monolith/internal/metrics"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store := metrics.NewEventStore()
	metricsSrv := &http.Server{Addr: ":8081", Handler: metrics.NewRouter(store)}
	iamSrv := &http.Server{Addr: ":8082", Handler: iampkg.NewMux()}

	go func() { log.Println(metricsSrv.ListenAndServe()) }()
	go func() { log.Println(iamSrv.ListenAndServe()) }()

	<-ctx.Done()
	stop()
	metricsSrv.Shutdown(context.Background()) //nolint:errcheck
	iamSrv.Shutdown(context.Background())     //nolint:errcheck
}
