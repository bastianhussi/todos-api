package profile

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
)

// fromRequest extracts the profile from the users request. NOTE: this does not check if any fields
// are provided. Default / nil values are allowed.
func fromRequest(r *http.Request) (*api.Profile, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	p := new(api.Profile)

	// Deserialize the json data. NOTE: Values can be nil / have the default type value.
	if err := json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}

// FIXME: Too many arguments! Maximum number of arguments should be 3.
func updateProfileInDB(ctx context.Context, conn *pg.Conn, id int, p *api.Profile) error {
	// TODO: is this really necessary?
	columns := make([]string, 3)
	if p.Email != "" {
		columns = append(columns, "email")
	}
	if p.Name != "" {
		columns = append(columns, "name")
	}
	if p.Password != "" {
		columns = append(columns, "password")
	}

	_, err := pg.Model(p).Limit(1).Column(columns...).Where("id = ?", id).Update()
	return err
}

func getProfileIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 0, 64)
	if err != nil {
		return 0, errors.New("Please profile a positive numeric value for the profile id")
	}

	return int(id), nil
}

// TODO: Add these functions as methods to the profile struct
func getProfileFromDB(ctx context.Context, conn *pg.Conn, id int) (*api.Profile, error) {
	profile := new(api.Profile)
	if err := conn.ModelContext(ctx, profile).Limit(1).Where("id = ?", id).Select(); err != nil {
		return nil, err
	}

	return profile, nil
}

func deleteProfileFromDB(ctx context.Context, conn *pg.Conn, id int) error {
	profile := new(api.Profile)
	if _, err := conn.ModelContext(ctx, profile).Limit(1).Where("id = ?", id).Delete(); err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}
