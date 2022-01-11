package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Service represents an external service
type Service interface {
	Do(action string)
}

// ServiceFactory is a component capable of constructing
// a service interface with a given service name
type ServiceFactory interface {
	New(serviceName string) (Service, error)
}

// stub service implementations
type FooService struct {
	Service
}

type BarService struct {
	Service
}

// default ServiceFactory implementation
type serviceFactory struct {
	// data
}

func (fac *serviceFactory) New(serviceName string) (Service, error) {
	serviceName = strings.TrimSpace(serviceName)
	switch serviceName {
	case "foo":
		return &FooService{ /* service specific config */ }, nil
	case "bar":
		return &BarService{ /* service specific config */ }, nil
	}
	return nil, fmt.Errorf("don't know how to construct %s", serviceName)
}

type serviceHandler struct {
	srvFactory ServiceFactory
}

func (handler *serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// the code below extracts the service name
	// and action name
	pathComponents := strings.Split(r.URL.Path, "/")
	l := len(pathComponents) - 1
	if pathComponents[l] == "/" {
		pathComponents = pathComponents[:l]
		l = len(pathComponents)
	}
	actionName := pathComponents[l]
	srvName := pathComponents[l-1]

	// construct the service from the factory
	// and call the specified action
	srv, err := handler.srvFactory.New(srvName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.Do(actionName)
}

func main() {
	handler := &serviceHandler{
		srvFactory: new(serviceFactory),
	}

	// requests will be of the form
	// /v0/dispatch/<serviceName>/<action>
	http.Handle("/v0/dispatch/", handler)
	http.ListenAndServe(":8080", nil)
}

// default ServiceFactory implementation
type serviceFactory1 struct {
	mu    sync.Mutex
	cache map[string]Service
}

func (fac *serviceFactory1) New1(serviceName string) (Service, error) {
	// check if the service has already been constructed -- (1)
	serviceName = strings.TrimSpace(serviceName)
	if srv, exists := fac.cache[serviceName]; exists {
		return srv, nil
	}

	fac.mu.Lock()
	defer fac.mu.Unlock()

	// ensure that the service wasn't already created during
	// lock-acquisition by a similar goroutine -- (2)
	if srv, exists := fac.cache[serviceName]; exists {
		return srv, nil
	}

	// construct the service and add it to the cache
	var newService Service
	switch serviceName {
	case "foo":
		newService = &FooService{ /* service specific config */ }
	case "bar":
		newService = &BarService{ /* service specific config */ }
	}
	if newService == nil {
		return nil, fmt.Errorf("don't know how to construct %s", serviceName)
	}
	fac.cache[serviceName] = newService
	return newService, nil
}