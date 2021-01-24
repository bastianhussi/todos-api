package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
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

	if profile.Password == "" {
		return http.StatusBadRequest, errors.New("Please provide a password")
	}

	c := h.res.DB.Conn()
	defer c.Close()

	storedProfile := new(api.Profile)

	if err := c.Model(storedProfile).Where("email = ?", profile.Email).Select(); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, errors.New(fmt.Sprintf("Could not find a profile with the email address %s", profile.Email))
		} else {
			panic(err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedProfile.Password), []byte(profile.Password))
	must(err)

	res, err := json.Marshal(storedProfile)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(res)
	must(err)

	return 0, nil
}

// Login handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer h.res.HandleRequestPanic(w, r)

	switch r.Method {
	case http.MethodPost:
		if code, err := h.post(w, r); err != nil {
			h.res.HandleRequest(w, r, code, err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Route add the routes of this package to the mux
func (h *Handler) Route(m *http.ServeMux) {
	m.HandleFunc("/login", h.res.Logging(h.Login))
	m.HandleFunc("/login/", h.res.Logging(h.Login))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
