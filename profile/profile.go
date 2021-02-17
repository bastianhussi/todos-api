package profile

import (
	"net/http"
	"strconv"

	api "github.com/bastianhussi/todos-api"
	"github.com/gorilla/mux"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

// FIXME: refactor this, so that these methods still implement http.Handler by removing the channel arugment.
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := api.DBFromContext(ctx)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, err := api.GetProfileByID(ctx, db, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.Respond(w, http.StatusOK, profile)
}

func (h *ProfileHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn := api.DBFromContext(ctx)

	reqProfile := new(api.Profile)
	if err := api.Decode(r, reqProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbProfile, err := api.GetProfileByID(ctx, conn, int(id))

	// FIXME: rollback if one of the tree transaction fails
	// FIXME: changing the email address should require authenticating the new email address.
	if err := dbProfile.Update(ctx, conn, reqProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	api.Respond(w, http.StatusOK, dbProfile)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn := api.DBFromContext(ctx)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, err := api.GetProfileByID(ctx, conn, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := profile.Delete(ctx, conn); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.Respond(w, http.StatusOK, profile)
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPatch:
		h.Patch(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	}
}
