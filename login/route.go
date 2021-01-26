package login

import (
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

var res *api.Resources

func post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	ctx := r.Context()

	p, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	conn := res.DB.Conn()
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

func login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	c := make(chan struct{})

	switch r.Method {
	case http.MethodPost:
		go post(w, r, c)

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

func NewHandler(s *api.Server) {
	res = s.Res
	s.AddRoute([]string{"/login"}, login, "POST")
}
