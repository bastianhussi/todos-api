package profile

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
)

func getProfileIDFromRequest(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)
	_, ok := vars["id"]
	if !ok {
		return 0, errors.New("Please provide a profile id")
	}

	id, err := strconv.ParseUint(vars["id"], 0, 64)
	if err != nil {
		return 0, errors.New("Please profile a positive numeric value for the profile id")
	}

	return id, nil
}

func getProfileFromDB(ctx context.Context, conn *pg.Conn, id uint64) (*api.Profile, error) {
	profile := new(api.Profile)
	if err := conn.ModelContext(ctx, profile).Limit(1).Where("id = ?", id).Select(); err != nil {
		return nil, err
	}

	return profile, nil
}

func deleteProfileFromDB(ctx context.Context, conn *pg.Conn, id uint64) error {
	profile := new(api.Profile)
	if _, err := conn.ModelContext(ctx, profile).Limit(1).Where("id = ?", id).Delete(); err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}
