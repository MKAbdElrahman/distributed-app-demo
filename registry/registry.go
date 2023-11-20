package registry

import (
	"errors"
	"sync"
)

type Registration struct {
	ServiceType          string   `json:"serviceType"`
	Port                 int      `json:"port"`
	IP                   string   `json:"ip"`
	RequiredServices     []string `json:"dependentServices"`
	NotificationEndpoint string   `json:"notificationEndpoint"`
}

type ServiceRegistry interface {
	GetServices() ([]Registration, error)
	GetDependentServices(serviceName string) ([]Registration, error)
	PostService(r *Registration) (*Registration, error)
	DeleteService(serviceName string) error
}

type InMemoryServiceRegistry struct {
	mu       sync.Mutex
	services []Registration
}

func (r *InMemoryServiceRegistry) GetServices() ([]Registration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.services, nil
}

func (r *InMemoryServiceRegistry) PostService(registration *Registration) (*Registration, error) {
	if registration == nil || registration.ServiceType == "" {
		return nil, errors.New("invalid registration")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existingService := range r.services {
		if existingService.ServiceType == registration.ServiceType &&
			existingService.IP == registration.IP &&
			existingService.Port == registration.Port {
			return nil, errors.New("service already registered with the same type, IP, and port")
		}
	}

	r.services = append(r.services, *registration)
	return registration, nil
}

func (r *InMemoryServiceRegistry) DeleteService(serviceName string) error {
	if serviceName == "" {
		return errors.New("invalid service name")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, existingService := range r.services {
		if existingService.ServiceType == serviceName {
			index = i
			break
		}
	}

	if index != -1 {
		r.services = append(r.services[:index], r.services[index+1:]...)
		return nil
	}

	return errors.New("service not found")
}

func (r *InMemoryServiceRegistry) GetDependentServices(serviceName string) ([]Registration, error) {
	if serviceName == "" {
		return nil, errors.New("invalid service name")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var dependentServices []Registration

	for _, existingService := range r.services {
		for _, dependentServiceName := range existingService.RequiredServices {
			if dependentServiceName == serviceName {
				dependentServices = append(dependentServices, existingService)
				break
			}
		}
	}

	return dependentServices, nil
}
