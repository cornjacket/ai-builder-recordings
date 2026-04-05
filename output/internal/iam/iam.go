// Package iam wires the lifecycle and authz sub-muxes into a single http.ServeMux.
//
// Purpose: Composes IAM sub-components into one mux returned by NewMux() for use by the platform main.
// Tags: implementation, iam
package iam

import (
	"net/http"
	"strings"

	"github.com/cornjacket/platform-monolith/internal/iam/internal/iam/authz"
	"github.com/cornjacket/platform-monolith/internal/iam/internal/iam/lifecycle"
)

// NewMux builds and returns a *http.ServeMux with all IAM routes registered.
// Lifecycle routes: /users, /auth/login, /auth/logout.
// Authz routes: /roles, /authz/check, /users/{id}/roles.
// /users/{id} (without /roles suffix) is forwarded to lifecycle.
func NewMux() *http.ServeMux {
	lc := lifecycle.Handler()
	az := authz.Handler()

	mux := http.NewServeMux()

	// lifecycle routes
	mux.Handle("/users", lc)
	mux.Handle("/auth/login", lc)
	mux.Handle("/auth/logout", lc)

	// authz routes
	mux.Handle("/roles", az)
	mux.Handle("/authz/check", az)

	// /users/{id} → lifecycle; /users/{id}/roles → authz
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/roles") {
			az.ServeHTTP(w, r)
		} else {
			lc.ServeHTTP(w, r)
		}
	})

	return mux
}
