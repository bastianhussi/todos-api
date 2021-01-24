package api

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var ctx = context.Background()

type Profile struct {
	tableName struct{} `pg:"profiles,alias:profile"`
	ID        int      `pg:",pk"`
	Email     string   `pg:",unique" json:"email"`
	Name      string   `pg:",notnull" json:"name"`
	Password  string   `pg:",notnull" json:"password"`
}

type Todo struct {
	tableName struct{}  `pg:"todos,alias:todo"`
	ID        int       `pg:",pk"`
	Title     string    `pg:",notnull"`
	CreatedAt time.Time `pg:"default:now()"`
	ProfileID int
	Profile   *Profile `pg:"rel:has-one,notnull"`
}

type Task struct {
	tableName struct{} `pg:"tasks,alias:task"`
	ID        int      `pg:",pk"`
	Title     string   `pg:",notnull"`
	Done      bool     `pg:"default:FALSE"`
	TodoID    int
	Todo      *Todo `pg:"rel:has-one,notnull"`
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
func CreateSchema(db *pg.Conn) error {
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
