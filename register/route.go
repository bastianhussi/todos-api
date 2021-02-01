package register

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

type Handler struct{}

// TODO: refactor this method. It should be stripped-down.
// post handles the incoming post request. When the request has been processed
// an empty struct is send into the channel indicating, that the task has been completed.
func (h *Handler) post(w http.ResponseWriter, r *http.Request, c chan<- struct{}) {
	panic("Ja moin")
	defer r.Body.Close()
	ctx := r.Context()

	conn := api.DBFromContext(ctx)

	profile := new(api.NewProfile)
	if err := api.Decode(r, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c <- struct{}{}
		return
	}

	dbProfile, err := profile.Insert(ctx, conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	api.Respond(w, r, http.StatusCreated, dbProfile)

	c <- struct{}{}
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
