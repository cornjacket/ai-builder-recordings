// Package metrics provides an in-memory event recording service for frontend user interactions.
//
// Purpose: HTTP router wiring for the metrics service, dispatching /events by method.
// Tags: implementation, metrics
package metrics

import "net/http"

// NewRouter returns an http.Handler with all metrics routes registered.
// The /events path dispatches to PostEvent or GetEvents based on the request method.
func NewRouter(store *EventStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostEvent(store)(w, r)
		case http.MethodGet:
			GetEvents(store)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
