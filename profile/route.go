package profile

import (
	"encoding/json"
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

type Handler struct {
	res *api.Resources
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(p)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := fromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	conn := h.res.DB.Conn()
	defer conn.Close()

	// FIXME: rollback if one of the tree transaction fails

	// FIXME: changing the email address should require authenticating the new email address.
	if err := updateProfileInDB(ctx, conn, p.ID, profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// FIXME: get updated profile directly from update queries
	profile, err = getProfileFromDB(ctx, conn, p.ID)
	must(err)

	body, err := json.Marshal(profile)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn := h.res.DB.Conn()
	defer conn.Close()

	// FIXME: is this always a bad request?
	if err := deleteProfileFromDB(ctx, conn, p.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.Respond(w, r, http.StatusOK, p)
}

func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

// TODO: use this to add it to the router, by implementing Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := getProfileIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn := h.res.DB.Conn()

	profile, err := getProfileFromDB(ctx, conn, id)
	if err != nil {
		http.Error(w, "Found so such user profile", http.StatusNotFound)
		return
	}

	// Close the connection already. Not doing so would cause two connections being open for each
	// request.
	conn.Close()
	// default not necessary: Mux already handles method not allowed cases.

}
