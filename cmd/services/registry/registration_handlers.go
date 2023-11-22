package main

import (
	"bytes"
	"demo/registry"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RegistrationHandler struct {
	Registry registry.ServiceRegistry
}

func (rh *RegistrationHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", rh.RegisterService)
	r.Get("/services", rh.GetServices)
	r.Delete("/deregister/{id}", rh.DeregisterService)
}

func (rh *RegistrationHandler) RegisterService(w http.ResponseWriter, r *http.Request) {
	var registration *registry.Registration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		log.Println("Failed to decode registration request:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	registeredService, err := rh.Registry.PostService(registration)
	if err != nil {
		log.Println("Failed to register service:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rh.findAndNotifyDependentServices("register", registration)
	if err != nil {
		log.Println("Failed to notify dependent services:", err)
		http.Error(w, "failed to notify dependent services", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registeredService)
}

func (rh *RegistrationHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := rh.Registry.GetServices()
	if err != nil {
		log.Println("failed to get services")
		http.Error(w, "failed to get services", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (rh *RegistrationHandler) DeregisterService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "id")

	if serviceID == "" {
		log.Println("Invalid service ID:", serviceID)
		http.Error(w, "invalid service ID", http.StatusBadRequest)
		return
	}

	service, err := rh.Registry.GetServiceByID(serviceID)

	if err != nil {
		log.Println("Failed to get service by ID:", err)
		http.Error(w, "failed to get service by ID", http.StatusInternalServerError)
		return
	}

	err = rh.findAndNotifyDependentServices("deregister", service)

	if err != nil {
		log.Println("Failed to notify dependent services:", err)
		http.Error(w, "failed to notify dependent services", http.StatusInternalServerError)
		return
	}

	err = rh.Registry.DeleteService(serviceID)
	if err != nil {
		log.Println("Failed to delete service:", err)
		http.Error(w, "failed to delete service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service deregistered successfully"))
}

func (rh *RegistrationHandler) findAndNotifyDependentServices(action string, service *registry.Registration) error {
	dependentServices, err := rh.Registry.GetDependentServices(service.ServiceType)
	if err != nil {
		return err
	}

	for _, dependentService := range dependentServices {
		notificationURL := dependentService.NotificationEndpoint

		payload := registry.NotificationPayload{
			Action:       action,
			Registration: *service,
		}
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		_, err = http.Post(notificationURL, "application/json", bytes.NewBuffer(payloadJSON))
		if err != nil {
			return err
		}
	}

	return nil
}
