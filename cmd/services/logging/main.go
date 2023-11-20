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

func setupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handler, err := DefaultLogToFileHandler("app.log")
	if err != nil {
		log.Fatal(err)
	}
	handler.RegisterRoutes(router)
	return router
}

func main() {
	port := flag.Int("port", 8081, "Port for the HTTP server")
	registrationAddr := flag.String("registration-addr", "http://localhost:8080/register", "Registration service endpoint")
	deregistrationAddr := flag.String("deregistration-addr", "http://localhost:8080/deregister", "Deregistration service endpoint")
	flag.Parse()

	server := &server.Server{
		Router:               setupRouter(),
		RegistrationAddr:     *registrationAddr,
		DeregistrationAddr:   *deregistrationAddr,
		Port:                 *port,
		ServiceName:          "Logging",
		RequiredServices:     []string{},
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
