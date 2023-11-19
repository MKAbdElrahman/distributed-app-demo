package registry

import (
	"errors"
	"sync"
)

type Registration struct {
	ServiceName string `json:"serviceName"`
	Port        int    `json:"port"`
	IP          string `json:"ip"`
}

type ServiceRegistry interface {
	GetServices() ([]Registration, error)
	PostService(r *Registration) (*Registration, error)
	DeleteService(serviceName string) error
}

type InMemoryServiceRegistry struct {
	mu       sync.Mutex
	services []Registration
}

// GetServices retrieves the list of registered services from in-memory storage.
func (r *InMemoryServiceRegistry) GetServices() ([]Registration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.services, nil
}

// PostService adds a new service registration to the in-memory storage.
func (r *InMemoryServiceRegistry) PostService(registration *Registration) (*Registration, error) {
	if registration == nil || registration.ServiceName == "" {
		return nil, errors.New("invalid registration")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the service is already registered
	for _, existingService := range r.services {
		if existingService.ServiceName == registration.ServiceName {
			return nil, errors.New("service already registered")
		}
	}

	// Add the new service registration
	r.services = append(r.services, *registration)
	return registration, nil
}

// DeleteService removes a service registration from in-memory storage.
func (r *InMemoryServiceRegistry) DeleteService(serviceName string) error {
	if serviceName == "" {
		return errors.New("invalid service name")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Find the index of the service with the given name
	index := -1
	for i, existingService := range r.services {
		if existingService.ServiceName == serviceName {
			index = i
			break
		}
	}

	// If the service is found, remove it from the slice
	if index != -1 {
		r.services = append(r.services[:index], r.services[index+1:]...)
		return nil
	}

	return errors.New("service not found")
}
