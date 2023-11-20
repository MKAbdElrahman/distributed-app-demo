package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct{}

type NotificationPayload struct {
	Action      string `json:"action"`
	ServiceName string `json:"serviceName"`
}

func (bh *NotificationHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/notify", bh.HandleReceivedNotification)
}

func (nh *NotificationHandler) HandleReceivedNotification(w http.ResponseWriter, r *http.Request) {
	var payload NotificationPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "failed to decode notification payload", http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "register":
		fmt.Printf("Received registration notification - Service: %s\n", payload.ServiceName)
	case "deregister":
		fmt.Printf("Received deregistration notification - Service: %s\n", payload.ServiceName)
	default:
		http.Error(w, "unknown action in notification", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification handled successfully"))
}
