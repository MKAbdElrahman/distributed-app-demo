package main

import (
	"demo/registry"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegistrationHandler is a struct that uses a ServiceRegistry to handle registrations.
type RegistrationHandler struct {
	Registry registry.ServiceRegistry
}

// RegisterRoutes sets up the HTTP routes for handling service registrations.
func (rh *RegistrationHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", rh.RegisterService)
	r.Get("/services", rh.GetServices)
	r.Delete("/deregister/{serviceName}", rh.DeregisterService)
}

// RegisterService handles the registration of a new service.
func (rh *RegistrationHandler) RegisterService(w http.ResponseWriter, r *http.Request) {
	var registration registry.Registration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	registeredService, err := rh.Registry.PostService(&registration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registeredService)
}

// GetServices handles the retrieval of registered services.
func (rh *RegistrationHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := rh.Registry.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// DeregisterService handles the deregistration of a service.
func (rh *RegistrationHandler) DeregisterService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "serviceName")
	if serviceName == "" {
		http.Error(w, "invalid service name", http.StatusBadRequest)
		return
	}

	err := rh.Registry.DeleteService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service deregistered successfully"))
}
