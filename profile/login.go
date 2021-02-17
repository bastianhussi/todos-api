package profile

import (
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"golang.org/x/crypto/bcrypt"
)

type (
	LoginHandler struct {
		sharedKey []byte
	}

	LoginProfile struct {
		profile *api.Profile
		token   string
	}
)

// FIXME: add the sharedkey to the context for routes like this
// NewHandler creates a hanlder for the login route
func NewLoginHandler(k []byte) *LoginHandler {
	return &LoginHandler{k}
}

func (p *LoginProfile) Public() interface{} {
	return map[string]interface{}{
		"profile": p.profile.Public(),
		"token":   p.token,
	}
}

// ServeHTTP handles the incoming requests for this route. In the case of the /login route,
// there are only post request, trying to log the client and receive a jwt token to authenticate
// themselfes with.
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginProfile := new(api.LoginProfile)
	if err := api.Decode(r, loginProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := api.DBFromContext(ctx)
	profile, err := loginProfile.Select(ctx, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: execute both in parallel
	token, err := api.GenerateJWT(h.sharedKey, loginProfile.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(loginProfile.Password)) != nil {
		http.Error(w, "Wrong password! Please try again", http.StatusBadRequest)
		return
	}

	api.Respond(w, http.StatusCreated, &LoginProfile{profile, token})
}
