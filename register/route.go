package register

import (
	api "github.com/bastianhussi/todos-api"
	"net/http"
)

type Handler struct {
	res *api.Resources
}

// NewHandler creates a new Handler.
func NewHandler(res *api.Resources) *Handler {
	return &Handler{
		res,
	}
}

// Register handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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
	m.HandleFunc("/register", h.res.Logging(h.Register))
	m.HandleFunc("/register/", h.res.Logging(h.Register))
}
