package register

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
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

func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	ctx := r.Context()

	p, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	p.Password = string(encryptedPass)

	conn := h.res.DB.Conn()
	defer conn.Close()

	dbChannel := make(chan error)
	go saveUserInDB(ctx, conn, p, dbChannel)

	select {
	case err = <-dbChannel:
		if err != nil {
			pgErr, ok := err.(pg.Error)
			if ok && pgErr.IntegrityViolation() {
				http.Error(w, "Account already exists", http.StatusBadRequest)
				c <- struct{}{}
				return
			}
			h.res.Logger.Println(err.Error())
		}
	case <-ctx.Done():
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte{})

	c <- struct{}{}
}

// Register handles the request for the `/login` route.
// Only POST-request are allowed.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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
			panic("Request canceled by client")
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
