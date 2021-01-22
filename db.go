package api

import (
	"context"

	"github.com/go-pg/pg/v10"
)

var ctx = context.Background()

func NewDB() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "user",
		Password: "pass",
		Database: "db_name",
	})
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
