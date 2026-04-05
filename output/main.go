// Package main is the entry point for the user-service HTTP server.
//
// Purpose: Wires the in-memory store and HTTP handlers together and starts the server on :8080.
// Tags: implementation, main
package main

import (
	"net/http"

	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/handlers"
	"github.com/cornjacket/ai-builder/acceptance-spec/sandbox/regressions/user-service/internal/userservice/store"
)

func main() {
	s := store.New()
	h := handlers.New(s)
	http.ListenAndServe(":8080", h)
}
