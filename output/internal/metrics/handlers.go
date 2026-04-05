// Package metrics provides an in-memory event recording service for frontend user interactions.
//
// Purpose: HTTP handler constructors for recording and retrieving user interaction events.
// Tags: implementation, metrics
package metrics

import (
	"encoding/json"
	"net/http"
)

var validEventTypes = map[string]bool{
	"click-mouse": true,
	"submit-form": true,
}

// PostEvent returns an http.HandlerFunc that decodes a JSON event from the request body,
// validates it, stores it, and responds with 201 and the stored event.
func PostEvent(store *EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Type    string                 `json:"type"`
			UserID  string                 `json:"userId"`
			Payload map[string]interface{} `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		if !validEventTypes[body.Type] {
			http.Error(w, "invalid event type", http.StatusBadRequest)
			return
		}

		stored := store.Add(Event{
			Type:    body.Type,
			UserID:  body.UserID,
			Payload: body.Payload,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(stored)
	}
}

// GetEvents returns an http.HandlerFunc that responds with the full list of stored events as JSON.
func GetEvents(store *EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := store.List()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(events)
	}
}
