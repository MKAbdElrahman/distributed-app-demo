package server

import (
	"bytes"
	"context"
	"demo/registry"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

// Server represents a service that handles service registration, deregistration,
// and serves as a central point for managing dependencies and notifications.
type Server struct {
	// Router is the Chi router used to define HTTP routes.
	Router *chi.Mux

	// RegistrationAddr is the address used by the server to register itself to the registry service.
	RegistrationAddr string

	// DeregistrationAddr is the address used by the server to deregister itself from the registry service..
	DeregistrationAddr string

	// Port is the port on which the server is running.
	Port int

	// ServiceName is the unique name of the service.
	ServiceName string

	// RequiredServices is a list of service names that this service depends on.
	RequiredServices []string

	// NotificationEndpoint is the URL where the server receives notifications.
	NotificationEndpoint string
}

func (s *Server) StartServer() error {
	serverAddr := fmt.Sprintf(":%d", s.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: s.Router,
	}

	go func() {
		log.Printf("Starting server on port %d...", s.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	if err := s.RegisterMe(); err != nil {
		log.Printf("Error registering server: %v", err)
	}

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Deregister before shutting down
	if err := s.DeregisterMe(); err != nil {
		log.Printf("Error deregistering server: %v", err)

	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped.")

	return nil
}

func (s *Server) RegisterMe() error {
	selfRegistration := registry.Registration{
		ServiceName:          s.ServiceName,
		Port:                 s.Port,
		IP:                   "127.0.0.1",
		RequiredServices:     s.RequiredServices,
		NotificationEndpoint: s.NotificationEndpoint,
	}

	body, err := json.Marshal(selfRegistration)
	if err != nil {
		return err
	}

	// Make a POST request to register itself
	resp, err := http.Post(s.RegistrationAddr, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *Server) DeregisterMe() error {

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/%v", s.DeregistrationAddr, s.ServiceName), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code during deregistration: %d", resp.StatusCode)
	}

	return nil
}
