package api

import (
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

func VerifyJWT(raw string, sharedKey []byte, audience []string) bool {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return false
	}

	cl := jwt.Claims{}
	if err := token.Claims(sharedKey, &cl); err != nil {
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

func GenerateJWT(sharedKey []byte, audience ...string) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: sharedKey},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	claims := jwt.Claims{
		Subject:   "subject",
		Issuer:    "issuer",
		Audience:  audience,
		NotBefore: jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
	}
	token, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}
