package register

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

// Handler holds all information needed to handle the incoming request (db-access, logging, eg.)
type Handler struct {
	res *api.Resources
}

// NewHandler creates a new Handler.
func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) error {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	profile := new(api.Profile)

	if err := json.Unmarshal(data, profile); err != nil {
		return err
	}

	if profile.Email == "" {
		return errors.New("Please provide an email address")
	}

	if profile.Name == "" {
		return errors.New("Please provide a profile name")
	}

	if profile.Password == "" {
		return errors.New("Please provide a password")
	}

	c := h.res.DB.Conn()
	defer c.Close()

	if _, err := c.Model(profile).Insert(); err != nil {
		return err
	}

	if err := c.Model(profile).Where("email = ?", profile.Email).Select(); err != nil {
		return err
	}


	res, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

	return nil
}

// Register handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer h.res.HandleInternalServerError(w, r)

	switch r.Method {
	case http.MethodPost:
		err := h.post(w, r)
		if err != nil {
			h.res.HandleBadRequest(w, r, err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Route add the routes of this package to the mux
func (h *Handler) Route(m *http.ServeMux) {
	m.HandleFunc("/register", h.res.Logging(h.Register))
	m.HandleFunc("/register/", h.res.Logging(h.Register))
}
