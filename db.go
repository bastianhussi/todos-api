package api

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var ctx = context.Background()

type User struct {
	Id       int    `json:"id" pg:",pk"`
	Email    string `json:"email" pg:",unique"`
	Name     string `json:"name pg:",notnull"`
	Password string `json:"password" pg:",notnull"`
}

type Todo struct {
	Id    int    `json:"id" pg:",pk"`
	Title string `json:title pg:",notnull"`
	User  *User  `json:user pg:"on_delete:CASCADE"`
}

type Task struct {
	Id    int    `json:id pg:",pk"`
	Title string `json:title pg:",notnull"`
	Done  bool   `json:done pg:"default:FALSE"`
	Todo  *Todo  `json:todo pg:"on_delete:CASCADE"`
}

// NewDB establishes a connection to the database and returns the database handle.
func NewDB(c *Config) (*pg.DB, error) {
	db := pg.Connect(c.Postgres)

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateSchema creates the corresponding database tables to the defined structs
// like User, Todo, eg.
func CreateSchema(c *pg.Conn) error {
	defer c.Close()

	models := []interface{}{
		(*User)(nil),
		(*Todo)(nil),
		(*Task)(nil),
	}

	for _, model := range models {
		err := c.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
