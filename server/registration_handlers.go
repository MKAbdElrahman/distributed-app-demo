package server

import (
	"demo/registry"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) RegisterNotifyRoute() {
	s.Router.Post("/notify", s.HandleReceivedNotification)
}

func (nh *Server) HandleReceivedNotification(w http.ResponseWriter, r *http.Request) {
	var payload registry.NotificationPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "failed to decode notification payload", http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "register":
		fmt.Printf("Received registration notification - Service: %s\n", payload.Registration.ServiceType)

		// Check if any of the required services is not attached to any instance
		for _, requiredService := range nh.RequiredServices {
			if _, exists := nh.ConnectedInstances[requiredService]; !exists {
				// Service is not attached to any instance, so add it to the ConnectedInstances
				instance := registry.ConnectedInstance{
					ID:   payload.Registration.ID,
					IP:   payload.Registration.IP,
					Port: payload.Registration.Port,
				}
				nh.ConnectedInstances[requiredService] = instance

				fmt.Printf("Added new instance for required service: %s\n", requiredService)
			}
		}

	case "deregister":
		fmt.Printf("Received deregistration notification - Service: %s\n", payload.Registration.ServiceType)

		// Check if the service is connected and if the IP and port match
		if instance, exists := nh.ConnectedInstances[payload.Registration.ServiceType]; exists {
			if instance.IP == payload.Registration.IP && instance.Port == payload.Registration.Port {
				// Deattach the service by removing it from ConnectedInstances
				delete(nh.ConnectedInstances, payload.Registration.ServiceType)

				fmt.Printf("Deregistered service: %s\n", payload.Registration.ServiceType)
			} else {
				fmt.Printf("Mismatch in IP or port for deregistering service: %s\n", payload.Registration.ServiceType)
			}
		} else {
			fmt.Printf("Service not found for deregistration: %s\n", payload.Registration.ServiceType)
		}
	default:
		http.Error(w, "unknown action in notification", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification handled successfully"))
}
