package api

import (
	"log"
	"os"

	"github.com/go-pg/pg/v10"
)

type Resources struct {
	Logger    *log.Logger
	DB        *pg.DB
	SharedKey string
}

func NewResources(c *Config) (*Resources, error) {
	db, err := NewDB(c)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)

	return &Resources{
		logger, db, c.JWTSecret,
	}, nil
}
