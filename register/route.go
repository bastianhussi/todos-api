package register

import (
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

// Handler holds all information needed to handle the incoming request (db-access, logging, eg.)
type Handler struct {
	res *api.Resources
}

// NewHandler creates a new Handler.
func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

// post handles the incoming post request. When the request has been processed
// an empty struct is send into the channel indicating, that the task has been completed.
func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	ctx := r.Context()

	p, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	// start the encrption of the user password in separate goroutine.
	passChannel := make(chan string)
	go encryptPassword(p.Password, passChannel)

	// create a new database connection
	// FIXME: can this operation block if a lot of conns are open? Use a goroutine instead?
	conn := h.res.DB.Conn()
	defer conn.Close()
	dbChannel := make(chan dbResult)

	// receive the encrypted password before writing it to the database.
	p.Password = <-passChannel
	go saveUserInDB(ctx, conn, p, dbChannel)

	// check if the write to the database could finish before the request was cancelt.
	select {
	case res := <-dbChannel:
		// handle the database response: Commit if the write as successful
		// or 
		err = func(tx *pg.Tx, err error) error {
			// Could not commit: The transaction was already rolled back.
			if err != nil {
				pgErr, ok := err.(pg.Error)
				if ok && pgErr.IntegrityViolation() {
					return fmt.Errorf("Profile with email %s already exists", p.Email)
				}
				panic(err)
			} else {
				defer tx.Close()
				if err := tx.Commit(); err != nil {
					panic(err)
				}
			}
			return nil
		}(res.tx, res.err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			c <- struct{}{}
			return
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
	// recover from potential panics (aka internal server errors)
	defer h.res.HandleRequestPanic(w)

	c := make(chan struct{})
	switch r.Method {
	case http.MethodPost:
		// process the post request in a separate goroutine.
		go h.post(w, r, c)

		// check if the goroutine finishes before the request is beeing canceld.
		select {
		case <-c:
			// request successfully handled
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

// little helperfunction which causes a panic if the error is not nil.
// NOTE: This should only be used for functions that can recover from panics.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
