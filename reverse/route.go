// Route the request to a location
package reverse

import (
	"net/http"
)

// Router matches incoming request to a specific location
type Router interface {
	// if error is not nil, the request wll be aborted and error will be proxied to client.
	// if location is nil and error is nil, that means that router did not find any matching location
	Route(req *http.Request) (Endpoint, error)
}

// Helper router that always the same location
type ConstRouter struct {
	endpoint Endpoint
}

func (m *ConstRouter) Route(req *http.Request) (Endpoint, error) {
	return m.endpoint, nil
}

type StdRouter struct{
}

func (r *StdRouter) Route(req *http.Request) (Endpoint, error) {
	return NewHttpEndpoint(req.URL)
}

var DefaultRouter Router = &StdRouter{}
