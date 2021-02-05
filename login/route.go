package login

import (
	"net/http"

	"github.com/bastianhussi/todos-api"
)

type Handler struct {
	sharedKey string
}

func NewHandler(k string) *Handler {
	return &Handler{k}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile := new(api.Profile)
	if err := api.Decode(r, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := api.DBFromContext(ctx)
	dbProfile, err := api.GetProfileByEmail(ctx, db, profile.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: execute both in parallel

	token, err := api.GenerateJWT(h.sharedKey, profile.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok := decryptPass(ctx, dbProfile.Password, profile.Password); !ok {
		http.Error(w, "Wrong password! Please try again", http.StatusBadRequest)
	}

	api.Respond(w, http.StatusCreated, token)
}
