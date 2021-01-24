package login

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

func receiveUserFromDB(ctx context.Context, conn *pg.Conn, email string) (*api.Profile, error) {
	tx, err := conn.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Close()

	p := new(api.Profile)
	if err := conn.ModelContext(ctx, p).Where("email = ?", email).Select(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// Commit on success.
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return p, nil
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
