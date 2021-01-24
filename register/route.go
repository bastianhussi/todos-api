package register

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"golang.org/x/crypto/bcrypt"
)

// Handler holds all information needed to handle the incoming request (db-access, logging, eg.)
type Handler struct {
	res *api.Resources
}

// NewHandler creates a new Handler.
func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) (int, error) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	profile := new(api.Profile)

	if err := json.Unmarshal(data, profile); err != nil {
		return http.StatusBadRequest, err
	}

	if profile.Email == "" {
		return http.StatusBadRequest, errors.New("Please provide an email address")
	}

	if profile.Name == "" {
		return http.StatusBadRequest, errors.New("Please provide a profile name")
	}

	if profile.Password == "" {
		return http.StatusBadRequest, errors.New("Please provide a password")
	}

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	c := h.res.DB.Conn()
	defer c.Close()

	profile.Password = string(encryptedPass)
	if _, err := c.Model(profile).Insert(); err != nil {
		panic(err)
	}

	// FIXME: remove this. No need to return the created profile to the user
	if err := c.Model(profile).Where("email = ?", profile.Email).Select(); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := json.Marshal(profile)
	if err != nil {
		must(err)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

	return http.StatusCreated, nil
}

// Register handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer h.res.HandleRequestPanic(w, r)

	switch r.Method {
	case http.MethodPost:
		code, err := h.post(w, r)
		h.res.HandleRequest(w, r, code, err)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Route add the routes of this package to the mux
func (h *Handler) Route(m *http.ServeMux) {
	m.HandleFunc("/register", h.res.Logging(h.Register))
	m.HandleFunc("/register/", h.res.Logging(h.Register))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
