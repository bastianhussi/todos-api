package register

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
)

func saveUserInDB(ctx context.Context, conn *pg.Conn, p *api.Profile, c chan<- error) {
	tx, err := conn.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Close()

	time.Sleep(10 * time.Second)

	if _, err := conn.ModelContext(ctx, p).Insert(); err != nil {
		_ = tx.Rollback()
		c <- err
		return
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	c <- nil
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

	if p.Name == "" {
		return nil, errors.New("Please provide a profile name")
	}

	if p.Password == "" {
		return nil, errors.New("Please provide a password")
	}

	return p, nil
}
