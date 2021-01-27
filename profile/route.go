package profile

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/gorilla/mux"
)

type Handler struct {
	res *api.Resources
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(res *api.Resources) *Handler {
	return &Handler{res}
}

func (h *Handler) Handle(mux *mux.Router) {
	mux.HandleFunc("/profile", h.profile).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
}
