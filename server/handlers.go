package server

import (
	"demo/registry"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationPayload struct {
	Action      string `json:"action"`
	ServiceType string `json:"serviceType"`
	Port        int    `json:"port"`
	IP          string `json:"ip"`
}

func (s *Server) RegisterNotifyRoute() {
	s.Router.Post("/notify", s.HandleReceivedNotification)
}

func (nh *Server) HandleReceivedNotification(w http.ResponseWriter, r *http.Request) {
	var payload NotificationPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "failed to decode notification payload", http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "register":
		fmt.Printf("Received registration notification - Service: %s\n", payload.ServiceType)

		// Check if any of the required services is not attached to any instance
		for _, requiredService := range nh.RequiredServices {
			if _, exists := nh.ConnectedInstances[requiredService]; !exists {
				// Service is not attached to any instance, so add it to the ConnectedInstances
				instance := registry.ConnectedInstanceAddr{
					IP:   payload.IP,
					Port: payload.Port,
				}
				nh.ConnectedInstances[requiredService] = instance

				fmt.Printf("Added new instance for required service: %s\n", requiredService)
			}
		}

	case "deregister":
		fmt.Printf("Received deregistration notification - Service: %s\n", payload.ServiceType)

		// Check if the service is connected and if the IP and port match
		if instance, exists := nh.ConnectedInstances[payload.ServiceType]; exists {
			if instance.IP == payload.IP && instance.Port == payload.Port {
				// Deattach the service by removing it from ConnectedInstances
				delete(nh.ConnectedInstances, payload.ServiceType)

				fmt.Printf("Deregistered service: %s\n", payload.ServiceType)
			} else {
				fmt.Printf("Mismatch in IP or port for deregistering service: %s\n", payload.ServiceType)
			}
		} else {
			fmt.Printf("Service not found for deregistration: %s\n", payload.ServiceType)
		}

	default:
		http.Error(w, "unknown action in notification", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification handled successfully"))
}
