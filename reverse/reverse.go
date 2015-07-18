package reverse

import (
	"net/http"
	"github.com/h2object/oxy/utils"
)

type ReverseOption func(r *ReverseRouter) error


// ErrorHandler is a functional argument that sets error handler of the server
func ErrorHandler(h utils.ErrorHandler) ReverseOption {
	return func(r *ReverseRouter) error {
		r.errHandler = h
		return nil
	}
}

func Logger(l utils.Logger) ReverseOption {
	return func(r *ReverseRouter) error {
		r.log = l
		return nil
	}
}

func Route(route Router) ReverseOption {
	return func(r *ReverseRouter) error {
		r.route = route
		return nil
	}
}

type ReverseRouter struct{
	next       http.Handler
	route 	   Router
	errHandler utils.ErrorHandler
	log        utils.Logger
}

func New(next http.Handler, opts ...ReverseOption) (*ReverseRouter, error) {
	rr := &ReverseRouter{
		next: next,
	}
	for _, o := range opts {
		if err := o(rr); err != nil {
			return nil, err
		}
	}
	if rr.route == nil {
		rr.route = DefaultRouter
	}
	if rr.errHandler == nil {
		rr.errHandler = utils.DefaultHandler
	}
	if rr.log == nil {
		rr.log = utils.NullLogger
	}
	return rr, nil
}

func (r *ReverseRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	endpoint, err := r.route.Route(req)
	if err != nil {
		r.errHandler.ServeHTTP(w, req, err)
		return
	}
	// make shallow copy of request before chaning anything to avoid side effects
	newReq := *req
	newReq.URL = endpoint.GetUrl()
	r.next.ServeHTTP(w, &newReq)
}

