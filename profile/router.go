package profile

import (
	api "github.com/bastianhussi/todos-api"
	"net/http"
)

var res *api.Resources

func profile(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(s *api.Server) {
	res = s.Res
	s.AddRoute([]string{"/profile/{id}", "/p/{id}"}, profile, "GET", "POST", "PATCH", "DELETE")
}