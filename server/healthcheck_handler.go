package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthCheckResult struct {
	ServiceID   string `json:"service_id"`
	ServiceType string `json:"service_type"`
	Healthy     bool   `json:"healthy"`
	Message     string `json:"message,omitempty"`
}

func (s *Server) RegisterHealthcheckRoute() {
	s.Router.Get("/healthcheck", s.HandleHealthCheck)
}

func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {

	// let's assume your service is always healthy.
	result := HealthCheckResult{
		ServiceID:   s.ID,
		ServiceType: s.ServiceType,
		Healthy:     true,
		Message:     "Service is healthy",
	}

	// Encode result as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode health check result: %v", err), http.StatusInternalServerError)
		return
	}
}
