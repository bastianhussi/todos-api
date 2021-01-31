package profile

import (
	"encoding/json"
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

type Handler struct {
	res *api.Resources
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, p *api.Profile) {
	body, err := json.Marshal(p)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request, p *api.Profile) {
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

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, p *api.Profile) {
	ctx := r.Context()

	conn := h.res.DB.Conn()
	defer conn.Close()

	// FIXME: is this always a bad request?
	if err := deleteProfileFromDB(ctx, conn, p.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(p)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
		h.get(w, r, profile)
	case http.MethodPatch:
		h.patch(w, r, profile)
	case http.MethodDelete:
		h.delete(w, r, profile)
	}
	// default not necessary: Mux already handles method not allowed cases.
}

func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

func (h *Handler) RegisterRoute(s *api.Server) {
	s.AddHandler([]string{"/profiles/{id}", "/p/{id}"}, h.Profile, http.MethodGet, http.MethodPatch, http.MethodDelete)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
