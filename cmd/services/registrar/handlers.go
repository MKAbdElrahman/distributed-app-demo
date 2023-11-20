package main

import (
	"bytes"
	"demo/registry"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RegistrationHandler struct {
	Registry registry.ServiceRegistry
}

func (rh *RegistrationHandler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", rh.RegisterService)
	r.Get("/services", rh.GetServices)
	r.Delete("/deregister/{serviceName}/{ip}/{port}", rh.DeregisterService)
}

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

	err = rh.notifyDependentServices("register", registration.ServiceType, registeredService.Port, registration.IP)
	if err != nil {
		http.Error(w, "failed to notify dependent services", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registeredService)
}

func (rh *RegistrationHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	services, err := rh.Registry.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (rh *RegistrationHandler) DeregisterService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "serviceName")
	portParam := chi.URLParam(r, "port")
	ip := chi.URLParam(r, "ip")

	if serviceName == "" || ip == "" || portParam == "" {
		http.Error(w, "invalid service name, port, or IP", http.StatusBadRequest)
		return
	}

	port, err := strconv.Atoi(portParam)
	if err != nil {
		http.Error(w, "invalid port", http.StatusBadRequest)
		return
	}

	err = rh.notifyDependentServices("deregister", serviceName, port, ip)
	if err != nil {
		http.Error(w, "failed to notify dependent services", http.StatusInternalServerError)
		return
	}

	err = rh.Registry.DeleteService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service deregistered successfully"))
}

func (rh *RegistrationHandler) notifyDependentServices(action string, serviceType string, port int, ip string) error {
	dependentServices, err := rh.Registry.GetDependentServices(serviceType)
	if err != nil {
		return err
	}

	for _, dependentService := range dependentServices {
		notificationURL := dependentService.NotificationEndpoint

		payload := map[string]interface{}{
			"action":      action,
			"serviceType": serviceType,
			"port":        port,
			"ip":          ip,
		}

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		_, err = http.Post(notificationURL, "application/json", bytes.NewBuffer(payloadJSON))
		if err != nil {
			fmt.Println("Failed to notify dependent service:", err)
		}
	}

	return nil
}
