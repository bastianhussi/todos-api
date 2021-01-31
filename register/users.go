package register

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	api "github.com/bastianhussi/todos-api"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
)

type dbResult struct {
	tx  *pg.Tx
	err error
}

// TODO: return the inserted id!
// saveUserInDB tries to write the given user profile to the database.
// The transaction is not beeing committed in this function scope, instead the transaction handles
// is beeing send into the channel. If the insert already fails it will be rolled back and
// the error is send into the channel.
func saveUserInDB(ctx context.Context, conn *pg.Conn, p *api.Profile, c chan<- dbResult) {
	// Start the transaction.
	tx, err := conn.Begin()
	if err != nil {
		panic(err)
	}

	// Try to write into the database. If that fails rollback
	if _, err := conn.ModelContext(ctx, p).Returning("id").Insert(); err != nil {
		defer tx.Close()
		_ = tx.Rollback()
		c <- dbResult{nil, err}
		return
	}

	// Was able to write to the database. The transaction is not committed yet.
	c <- dbResult{tx, err}
}

// encryptPassword uses the bcrypt algorithm to encrypt a given password.
// The process of creating a hashed password can not be canceld using a context.
func encryptPassword(p string, c chan<- string) {
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	must(err)
	c <- string(encryptPass)
}
