package profile

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// FIXME: refactor this, so that these methods still implement http.Handler by removing the channel arugment.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, ch chan struct{}) {
	profile := new(api.Profile)
	if err := api.Decode(r, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// TODO: receive user from DB

	api.Respond(w, r, http.StatusOK, profile)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request, ch chan struct{}) {
	ctx := r.Context()
	id, err := getProfileIDFromRequest(r)
	if err != nil {
		respondWithBadRequest(w, err)
		return
	}

	profile, err := fromRequest(r)
	if err != nil {
		respondWithBadRequest(w, err)
		return
	}

	db, _ := ctx.Value("db").(*pg.Conn)

	// FIXME: rollback if one of the tree transaction fails

	// FIXME: changing the email address should require authenticating the new email address.
	if err := updateProfileInDB(ctx, db, id, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// FIXME: get updated profile directly from update queries
	profile, err = getProfileFromDB(ctx, db, id)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, ch chan struct{}) {

}

// TODO: Extract this function from this package and move it to the parent package.
func respondWithBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ch := make(chan struct{}, 1)

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, ch)
	case http.MethodPatch:
		h.Patch(w, r, ch)
	case http.MethodDelete:
		h.Delete(w, r, ch)
	}

	select {
	case <-ch:
		return
	case <-ctx.Done():
		return
	}
}
