package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Nice blog post about pg: https://medium.com/tunaiku-tech/go-pg-golang-postgre-orm-2618b75c0430

type Public interface {
	Public() interface{}
}

type NewProfile struct {
	Email           string
	Name            string
	Password        string
	PasswordConfirm string
}

func requiredError(v string) error {
	return errors.New(fmt.Sprintf("%s is required", v))
}

func (p *NewProfile) OK() error {
	if len(p.Name) == 0 {
		return requiredError("Name")
	}

	if len(p.Password) == 0 || len(p.PasswordConfirm) == 0 {
		return requiredError("Password")
	}

	if len(p.Password) < 8 || len(p.PasswordConfirm) < 8 {
		return errors.New("Password must contain at least 8 characters")
	}

	if p.Password != p.PasswordConfirm {
		return errors.New("Passwords don't match")
	}

	return nil
}

func (p *NewProfile) Save(conn *pg.Conn) (*Profile, error) {
	// TODO: implement

	return nil, nil
}

// Profile
type Profile struct {
	tableName    struct{} `pg:"profiles,alias:profile"`
	ID           int      `pg:",pk"`
	Email        string   `pg:",unique" json:"email"`
	Name         string   `pg:",notnull" json:"name"`
	PasswordHash string   `pg:",notnull" json:"-"`
}

func (p *Profile) Public() interface{} {
	return map[string]interface{}{
		"id":    p.ID,
		"email": p.Email,
		"name":  p.Name,
	}
}

func (p *Profile) OK() error {
	if len(p.Name) == 0 {
		return errors.New("Name is required")
	}

	if len(p.Email) == 0 {
		return errors.New("Email is required")
	}

	if len(p.Password) == 0 {
		return errors.New("Password is required")
	}

	return nil
}

// Todo
type Todo struct {
	tableName struct{}  `pg:"todos,alias:todo"`
	ID        int       `pg:",pk"`
	Title     string    `pg:",notnull"`
	CreatedAt time.Time `pg:"default:now()"`
	ProfileID int
	Profile   *Profile `pg:"rel:has-one,notnull"`
}

// Task
type Task struct {
	tableName struct{} `pg:"tasks,alias:task"`
	ID        int      `pg:",pk"`
	Title     string   `pg:",notnull"`
	Done      bool     `pg:"default:FALSE"`
	TodoID    int
	Todo      *Todo `pg:"rel:has-one,notnull"`
}

// NewDB establishes a connection to the database and returns the database handle.
func NewDB(ctx context.Context, c *Config) (*pg.DB, error) {
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
