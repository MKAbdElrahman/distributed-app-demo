package registry

import (
	"errors"
	"sync"
)

// ConnectedInstanceAddr represents the address information of a connected instance.
type ConnectedInstance struct {
	ID   string `json:"int"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// ConnectedInstance represents a specific instance of a connected service.
// a map from a service type to connected instance
type ConnectedInstances map[string]ConnectedInstance

type Registration struct {
	ID                   string             `json:"id"`
	ServiceType          string             `json:"serviceType"`
	Port                 int                `json:"port"`
	IP                   string             `json:"ip"`
	RequiredServices     []string           `json:"dependentServices"`
	ConnectedInstances   ConnectedInstances `json:"connectedInstances"`
	NotificationEndpoint string             `json:"notificationEndpoint"`
	HealthCheckEndpoint  string             `json:"healthcheckEndpoint"`
}

type ServiceRegistry interface {
	GetServices() ([]Registration, error)
	GetServicesByType(name string) ([]Registration, error)
	GetServiceByID(id string) (*Registration, error)
	GetDependentServices(serviceName string) ([]Registration, error)
	PostService(r *Registration) (*Registration, error)
	DeleteService(serviceName string) error
}

type InMemoryServiceRegistry struct {
	mu       sync.Mutex
	services []Registration
}

type NotificationPayload struct {
	Action       string `json:"action"`
	Registration Registration
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

func (r *InMemoryServiceRegistry) DeleteService(serviceID string) error {
	if serviceID == "" {
		return errors.New("invalid service ID")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, existingService := range r.services {
		if existingService.ID == serviceID {
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

func (r *InMemoryServiceRegistry) GetServicesByType(serviceType string) ([]Registration, error) {
	if serviceType == "" {
		return nil, errors.New("invalid service type")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var servicesByType []Registration

	for _, existingService := range r.services {
		if existingService.ServiceType == serviceType {
			servicesByType = append(servicesByType, existingService)
		}
	}

	return servicesByType, nil
}

func (r *InMemoryServiceRegistry) GetServiceByID(id string) (*Registration, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existingService := range r.services {
		if existingService.ID == id {
			return &existingService, nil
		}
	}

	return nil, errors.New("id not found")
}
