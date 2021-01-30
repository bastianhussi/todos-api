package api

import (
	"github.com/square/go-jose/v3/jwt"
)

func verifyToken(raw string, sharedKey string, audience []string) bool {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return false
	}

	cl := jwt.Claims{}
	if err := token.Claims([]byte(sharedKey), &cl); err != nil {
		return false
	}

	// TODO: check if the audience is correct.
	if err = cl.Validate(jwt.Expected{
		Issuer:  "issuer",
		Subject: "subject",
		// Audience: audience,
	}); err != nil {
		return false
	}

	return true
}
