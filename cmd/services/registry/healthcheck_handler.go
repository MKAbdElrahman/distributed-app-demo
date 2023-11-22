package main

import (
	"demo/registry"
	"demo/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthCheckHandler struct {
	Registry registry.ServiceRegistry
}

func (rh *HealthCheckHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/healthchecks", rh.HandleHealthCheck)
}

// HandleHealthCheck pings all registered services to ensure they are up and running
func (rh *HealthCheckHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	services, err := rh.Registry.GetServices()
	if err != nil {
		log.Println("Failed to get services for health check:", err)
		http.Error(w, "failed to get services for health check", http.StatusInternalServerError)
		return
	}

	var healthCheckResults []server.HealthCheckResult

	for _, service := range services {
		healthCheckURL := service.HealthCheckEndpoint

		resp, err := http.Get(healthCheckURL)
		if err != nil {
			result := server.HealthCheckResult{
				ServiceID:   service.ID,
				ServiceType: service.ServiceType,
				Healthy:     false,
				Message:     fmt.Sprintf("Health check failed: %v", err),
			}
			healthCheckResults = append(healthCheckResults, result)
			continue
		}
		defer resp.Body.Close()

		healthy := resp.StatusCode == http.StatusOK
		message := fmt.Sprintf("Health check %s", http.StatusText(resp.StatusCode))

		result := server.HealthCheckResult{
			ServiceID:   service.ID,
			ServiceType: service.ServiceType,
			Healthy:     healthy,
			Message:     message,
		}
		healthCheckResults = append(healthCheckResults, result)
	}

	// Encode results as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(healthCheckResults); err != nil {
		log.Println("Failed to encode health check results:", err)
		http.Error(w, "failed to encode health check results", http.StatusInternalServerError)
		return
	}
}
