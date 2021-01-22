package login

import (
	"log"
	"net/http"
	"time"
)

// Handler manages all the routes and middleware for the login functionality.
type Handler struct {
	logger *log.Logger
}

// NewHandler creates a new Handler struct with the given logger.
func NewHandler(l *log.Logger) *Handler {
	return &Handler{l}
}

// Logging is a middleware function which keeps track of the request method
// and the time needed to process this request.
func (h *Handler) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Printf("Got %s request.\n", r.Method)
		start := time.Now()
		defer next(w, r)
		h.logger.Printf("Request processed in %s\n", time.Now().Sub(start))
	}

}

// Login handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		defer w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		w.Header().Add("Content-Type", "plain/text; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Hello, World!"))
	}
}

// Route add the routes of this package to the mux
func (h *Handler) Route(m *http.ServeMux) {
	m.HandleFunc("/login", h.Logging(h.Login))
	m.HandleFunc("/login/", h.Logging(h.Login))
}
