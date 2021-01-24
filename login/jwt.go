package login

import (
	api "github.com/bastianhussi/todos-api"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

var sharedKey = []byte("secret")

type TokenResult struct {
	token string
	err   error
}

func generateToken(p *api.Profile, c chan<- TokenResult) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: sharedKey}, (&jose.SignerOptions{}).WithType("JWT"))

	if err != nil {
		c <- TokenResult{err: err}
		return
	}

	claims := jwt.Claims{
		Subject:  "subject",
		Issuer:   "issuer",
		Audience: jwt.Audience{p.Email},
	}
	token, err := jwt.Signed(sig).Claims(claims).CompactSerialize()

	c <- TokenResult{token, err}
}
