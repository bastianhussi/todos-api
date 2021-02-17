package api

import (
	"encoding/json"
	"net/http"
)

type Public interface {
	Public() interface{}
}

func Must(err interface{}) {
	errVal, ok := err.(error)
	if ok {
		if errVal != nil {
			panic(err)
		}
		return
	}

	boolVal, ok := err.(bool)
	if ok {
		if boolVal {
			panic(err)
		}
	}
}

func Respond(w http.ResponseWriter, status int, data interface{}) {
	if obj, ok := data.(Public); ok {
		data = obj.Public()
	}

	body, err := json.Marshal(data)
	Must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func Decode(r *http.Request, v interface{ OK() error }) error {
	// This can check if the OK method on a struct returns an error.
	// We can check if required fields are given this way.
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	if err := v.OK(); err != nil {
		return err
	}

	return nil
}
