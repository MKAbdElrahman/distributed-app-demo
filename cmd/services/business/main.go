package main

import (
	"demo/server"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter(logHandler LogHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handler := BusinessHandler{
		LogHandler: logHandler,
	}

	noifyHandler := NotificationHandler{}

	noifyHandler.RegisterRoutes(router)
	handler.RegisterRoutes(router)

	return router
}

func main() {
	port := flag.Int("port", 8082, "Port for the HTTP server")
	registrationAddr := flag.String("registration-addr", "http://localhost:8080/register", "Registration service endpoint")
	deregistrationAddr := flag.String("deregistration-addr", "http://localhost:8080/deregister", "Deregistration service endpoint")
	loggingServiceURL := flag.String("logging-service-url", "http://localhost:8081", "URL of the logging service")
	flag.Parse()

	logHandler := &HTTPLogHandler{
		LoggerURL: *loggingServiceURL,
	}

	server := &server.Server{
		Router:               setupRouter(logHandler),
		RegistrationAddr:     *registrationAddr,
		DeregistrationAddr:   *deregistrationAddr,
		Port:                 *port,
		ServiceName:          "Business",
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
