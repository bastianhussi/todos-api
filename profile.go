package api

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
)

// NewProfile is a user profile that has not been saved to the database. It maches the form data a
// client uploads the create a profile.
type NewProfile struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (p *NewProfile) OK() error {
	if len(p.Name) == 0 {
		return RequiredError("Name")
	}

	if len(p.Email) == 0 {
		return RequiredError("Email")
	}

	if len(p.Password) == 0 {
		return RequiredError("Password")
	}

	if len(p.PasswordConfirm) == 0 {
		return RequiredError("PasswordConfirm")
	}

	if p.Password != p.PasswordConfirm {
		return errors.New("Passwords don't match")
	}

	if len(p.Password) < 8 {
		return errors.New("Password must contain at least 8 characters")
	}

	return nil
}

func (p *NewProfile) Insert(ctx context.Context, conn *pg.Conn) (*Profile, error) {
	tx, err := conn.Begin()
	defer tx.Close()
	Must(err)

	encryptPass, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	Must(err)

	profile := new(Profile)
	profile.Email = p.Email
	profile.Name = p.Name
	profile.Password = string(encryptPass)

	_, err = conn.ModelContext(ctx, profile).Returning("id").Insert()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return profile, nil
}

type LoginProfile struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *LoginProfile) OK() error {
	if len(p.Email) == 0 {
		return RequiredError("Email")
	}

	if len(p.Password) == 0 {
		return RequiredError("Password")
	}

	return nil
}

// Select searches the database for a user profile matching the fields from the submitted login form.
func (p *LoginProfile) Select(ctx context.Context, conn *pg.Conn) (*Profile, error) {
	profile := new(Profile)
	err := conn.ModelContext(ctx, profile).Limit(1).Where("email = ?", p.Email).Select()
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Profile
type Profile struct {
	tableName struct{} `pg:"profiles,alias:profile"`
	ID        int      `pg:",pk" json:"id"`
	Email     string   `pg:",notnull,unique" json:"email"`
	Name      string   `pg:",notnull" json:"name"`
	Password  string   `pg:",notnull" json:"password"`
}

func (p *Profile) Public() interface{} {
	return map[string]interface{}{
		"id":    p.ID,
		"email": p.Email,
		"name":  p.Name,
	}
}

func (p *Profile) OK() error {
	if len(p.Email) == 0 {
		return RequiredError("Email")
	}

	if len(p.Name) == 0 {
		return RequiredError("Name")
	}

	if len(p.Password) == 0 {
		return RequiredError("Password")
	}

	return nil
}

func (p *Profile) Update(ctx context.Context, conn *pg.Conn, profile *Profile) error {
	// TODO: implement

	return nil
}

// Delete removes the user profile from the database
func (p *Profile) Delete(ctx context.Context, conn *pg.Conn) error {
	tx, err := conn.Begin()
	Must(err)
	defer tx.Close()

	_, err = conn.ModelContext(ctx, p).Delete()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// FIXME: Do not commit here already, return the transaction handler.
	tx.Commit()

	return nil
}

// GetProfileByID searches the database for a user profile with the given ID and returns the record
// if it exists.
func GetProfileByID(ctx context.Context, conn *pg.Conn, id int) (*Profile, error) {
	profile := new(Profile)
	err := conn.ModelContext(ctx, profile).Limit(1).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}

	return profile, nil
}
