package register

import (
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

type Handler struct {
	res *api.Resources
}

// TODO: refactor this method. It should be stripped-down.
// post handles the incoming post request. When the request has been processed
// an empty struct is send into the channel indicating, that the task has been completed.
func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	defer r.Body.Close()
	ctx := r.Context()

	p, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	// TODO: remove these goroutines: There is no benefit in running encrption and the databse insert
	// in goroutines. The database operation has to wait on the other tasks anyways.
	// Could use a channel to send the password into the database function but the password is
	// required very early on in the function.

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

	// TODO: extract this into an other function.
	// This method should not handle the cancellation of this task.
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
					return fmt.Errorf("profile with email %s already exists", p.Email)
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
	_, _ = fmt.Fprintf(w, "Profile succesfully created")

	c <- struct{}{}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c := make(chan struct{}, 1)
	go h.post(w, r, c)

	select {
	case <-c:
		return
	case <-ctx.Done():
		panic(ctx.Err().Error())
	}
}

func (h *Handler) RegisterRoute(s *api.Server) {
	s.AddHandler([]string{"/register"}, h.Register, "POST")
}

func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

// little helperfunction which causes a panic if the error is not nil.
// NOTE: This should only be used for functions that can recover from panics.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
