package main

import (
	"demo/cmd/services/business/handlers"
	"demo/server"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter(logger handlers.HTTPLogger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handler := handlers.LogHandler{
		Logger: logger,
	}

	notifyHandler := handlers.NotificationHandler{}

	notifyHandler.RegisterRoutes(router)
	handler.RegisterRoutes(router)

	return router
}

func main() {
	port := flag.Int("port", 8082, "Port for the HTTP server")
	registrationAddr := flag.String("registration-addr", "http://localhost:8080/register", "Registration service endpoint")
	deregistrationAddr := flag.String("deregistration-addr", "http://localhost:8080/deregister", "Deregistration service endpoint")
	loggingServiceURL := flag.String("logging-service-url", "http://localhost:8081", "URL of the logging service")
	flag.Parse()

	logger := handlers.HTTPLogger{
		Endpoint: *loggingServiceURL + "/log",
	}

	server := &server.Server{
		Router:               setupRouter(logger),
		RegistrationAddr:     *registrationAddr,
		DeregistrationAddr:   *deregistrationAddr,
		Port:                 *port,
		ServiceType:          "Business",
		RequiredServices:     []string{"Logging"},
		NotificationEndpoint: fmt.Sprintf("http://localhost:%d/notify", *port),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.StartServer(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	wg.Wait()
}
