package api

import (
	"log"
	"net/http"
	"time"
)

// Routers can add endpoints to the servers mux.
type Router interface {
	Route(mux *http.ServeMux)
}

// Resources are used by the Handlers to access databases, use logging utilities, ...
type Resources struct {
	logger *log.Logger
}

// NewHandler creates a new Handler.
func NewResources(l *log.Logger) *Resources {
	return &Resources{l}
}

// Login handles the request for the `/login` route.
// Only POST-request are allowed.
func (res *Resources) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res.logger.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
		start := time.Now()
		defer next(w, r)
		res.logger.Printf("Request processed in %s\n", time.Now().Sub(start))
	}
}
