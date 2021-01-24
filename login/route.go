package login

import (
	"fmt"
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

func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	ctx := r.Context()

	p, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	conn := h.res.DB.Conn()
	defer conn.Close()

	// TODO: use a goroutine instead
	storedProfile, err := receiveUserFromDB(ctx, conn, p.Email)
	if err != nil {
		if err == pg.ErrNoRows {
			http.Error(w, fmt.Sprintf("No profile with email %s found", p.Email), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		c <- struct{}{}
		return
	}

	tokenRes := make(chan TokenResult)
	go generateToken(p, tokenRes)

	res := <-tokenRes
	if res.err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	bcrypt.CompareHashAndPassword([]byte(storedProfile.Password), []byte(p.Password))

	w.Header().Add("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(res.token))

	c <- struct{}{}
}

// Login handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer h.res.HandleRequestPanic(w)

	c := make(chan struct{})

	switch r.Method {
	case http.MethodPost:
		go h.post(w, r, c)

		select {
		case <-c:
			return
		case <-ctx.Done():
			err := ctx.Err()
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
