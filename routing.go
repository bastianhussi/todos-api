package api

import (
	"net/http"
)

// Routers can add endpoints to the servers mux.
type Router interface {
	Route(mux *http.ServeMux)
}
