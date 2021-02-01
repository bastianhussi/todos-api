package api

import (
	"errors"
	"fmt"
)

func RequiredError(v string) error {
	return errors.New(fmt.Sprintf("%s is required", v))
}

type BadRequestError struct {
	msg string
	err error
}

// TODO: add method to write a response to the user with the error message and a status code.

func (e *BadRequestError) Error() string {
	return e.msg
}

func (e *BadRequestError) Unwrap() error {
	return e.err
}

type DBError struct {
	msg string
	err error
}

func NewDBError(msg string, err error) *DBError {
	return &DBError{msg, err}
}

func (e *DBError) Error() string {
	return e.msg
}

func (e *DBError) Unwrap() error {
	return e.err
}
