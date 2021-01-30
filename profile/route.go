package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/bastianhussi/todos-api"
)

type Handler struct {
	res *api.Resources
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, id uint64) {
	ctx := r.Context()
	conn := h.res.DB.Conn()
	defer conn.Close()

	profile, err := getProfileFromDB(ctx, conn, id)
	if err != nil {
		http.Error(w, "Found so such user profile", http.StatusNotFound)
		return
	}

	body, err := json.Marshal(profile)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request, id uint64) {
	fmt.Println(w, "Hello, world!")
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id uint64) {
	ctx := r.Context()

	conn := h.res.DB.Conn()
	defer conn.Close()

	// FIXME: is this always a bad request?
	if err := deleteProfileFromDB(ctx, conn, id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	// TODO: return deleted user profile
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	id, err := getProfileIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPatch:
		h.patch(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
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
