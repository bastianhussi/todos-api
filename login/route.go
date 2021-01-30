package login

import (
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

type Handler struct {
	res *api.Resources
}

func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

func (h *Handler) RegisterRoute(s *api.Server) {
	s.AddHandler([]string{"/login"}, h.Login, "POST")
}

// TODO: only create goroutines if two or more tasks can be run in parallel
// TODO: context cancellation with select-statements are only necessary in goroutines
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
	dbProfile, err := receiveUserFromDB(ctx, conn, p.Email)
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

	if ok := decryptPass(ctx, dbProfile.Password, p.Password); !ok {
		http.Error(w, "Wrong password! Please try again", http.StatusBadRequest)
	}

	w.Header().Add("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(res.token))

	c <- struct{}{}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
