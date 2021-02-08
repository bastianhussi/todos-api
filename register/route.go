package register

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := api.DBFromContext(ctx)

	profile := new(api.NewProfile)
	if err := api.Decode(r, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbProfile, err := profile.Insert(ctx, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.Respond(w, http.StatusCreated, dbProfile)
}
