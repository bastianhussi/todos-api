package login

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

func receiveUserFromDB(ctx context.Context, conn *pg.Conn, email string) (*api.Profile, error) {
	p := new(api.Profile)
	if err := conn.ModelContext(
		ctx,
		p,
	).Limit(1).Where("email = ?", email).Select(); err != nil {
		return nil, err
	}

	return p, nil
}

func decryptPass(ctx context.Context, hashedPass string, pass string) bool {
	c := make(chan error, 1)
	go func() {
		c <- bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	}()

	select {
	case err := <-c:
		return err == nil
	case <-ctx.Done():
		return false
	}
}

func fromRequest(r *http.Request) (*api.Profile, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	p := new(api.Profile)

	if err := json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	if p.Email == "" {
		return nil, errors.New("Please provide an email address")
	}

	if p.Password == "" {
		return nil, errors.New("Please provide a password")
	}

	return p, nil
}
