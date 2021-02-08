package api

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Nice blog post about pg: https://medium.com/tunaiku-tech/go-pg-golang-postgre-orm-2618b75c0430

// NewDB establishes a connection to the database and returns the database handle.
func newDB(ctx context.Context, c *Config) (*pg.DB, error) {
	db := pg.Connect(c.Postgres)

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateSchema creates the corresponding database tables to the defined structs
// like User, Todo, eg.
func createSchema(db *pg.Conn) error {
	models := []interface{}{
		(*Profile)(nil),
		(*Todo)(nil),
		(*Task)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
