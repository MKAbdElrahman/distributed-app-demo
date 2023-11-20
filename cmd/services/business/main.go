package main

import (
	"demo/cmd/services/business/handlers"
	"demo/registry"
	"demo/server"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := flag.Int("port", 8082, "Port for the HTTP server")
	registrationAddr := flag.String("registration-addr", "http://localhost:8080/register", "Registration service endpoint")
	deregistrationAddr := flag.String("deregistration-addr", "http://localhost:8080/deregister", "Deregistration service endpoint")
	flag.Parse()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	server := &server.Server{
		Router:               router,
		RegistrationAddr:     *registrationAddr,
		DeregistrationAddr:   *deregistrationAddr,
		Port:                 *port,
		ServiceType:          "Business",
		RequiredServices:     []string{"Logging"},
		ConnectedInstances:   make(registry.ConnectedInstances),
		NotificationEndpoint: fmt.Sprintf("http://localhost:%d/notify", *port),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Start a goroutine to dynamically update the loggingServiceURL
	go func() {
		defer wg.Done()

		// Create an HTTPLogger with an empty endpoint initially
		logger := handlers.HTTPLogger{}

		// Assume we are waiting for the Logging service to connect
		for {
			// Access the ConnectedInstances to get the URL of the Logging service
			if loggingInstance, exists := server.ConnectedInstances["Logging"]; exists {
				// Update the HTTPLogger with the correct endpoint
				logger.Endpoint = fmt.Sprintf("http://%s:%d/log", loggingInstance.IP, loggingInstance.Port)

				handler := handlers.LogHandler{
					Logger: logger,
				}

				handler.RegisterRoutes(router)

				break
			}

			log.Println("Waiting for the Logging service to connect...")

			// Sleep for a short duration before retrying
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.StartServer(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	wg.Wait()

}
