package login

import (
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

type Handler struct {
	sharedKey string
}

func NewHandler(k string) *Handler {
	return &Handler{k}
}

// NOTE: only create goroutines if two or more tasks can be run in parallel
// NOTE: context cancellation with select-statements are only necessary in goroutines
func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	ctx := r.Context()
	profile := new(api.Profile)
	if err := api.Decode(r, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	db, _ := ctx.Value("db").(*pg.Conn)

	// TODO: use a goroutine instead
	dbProfile, err := receiveUserFromDB(ctx, db, profile.Email)
	if err != nil {
		if err == pg.ErrNoRows {
			http.Error(w, fmt.Sprintf("No profile with email %s found", profile.Email), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		c <- struct{}{}
		return
	}

	token, err := api.GenerateJWT(h.sharedKey, profile.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	if ok := decryptPass(ctx, dbProfile.Password, profile.Password); !ok {
		http.Error(w, "Wrong password! Please try again", http.StatusBadRequest)
	}

	api.Respond(w, r, http.StatusCreated, token)

	c <- struct{}{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c := make(chan struct{})
	go h.post(w, r, c)

	select {
	case <-c:
		return
	case <-ctx.Done():
		err := ctx.Err()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
