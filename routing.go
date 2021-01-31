package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

func WithLogger(l *log.Logger, h http.Handler) http.Handler {
	return &logwrapper{l, h}
}

type (
	logwrapper struct {
		logger  *log.Logger
		handler http.Handler
	}
	dbwrapper struct {
		dbSession *pg.DB
		handler   http.Handler
	}
)

func WithDB(d *pg.DB, h http.Handler) http.Handler {
	return &dbwrapper{d, h}
}

func (l *logwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "logger", l.logger)
	r.WithContext(ctx)

	l.handler.ServeHTTP(w, r)
}

// Provide a open db connection for each request using this and make sure the connection is closed
// when finished
func (d *dbwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := d.dbSession.Conn()
	defer conn.Close()

	ctx := r.Context()
	ctx = context.WithValue(ctx, "db", conn)
	r.WithContext(ctx)

	d.handler.ServeHTTP(w, r)
}

type Adapter func(http.Handler) http.Handler

func Logging(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			l.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
			l.Printf("Request processed in %s\n", time.Since(start))
		})
	}
}

// TODO: does this work like this? (Using it as the outermost adapter)
func Recover(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l.Println(err)
					http.Error(w, "Oh noo, something went wrong 🤯", http.StatusInternalServerError)
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

func Auth() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.Split(token, "Bearer ")[1]

			if ok := VerifyJWT(token, "secret", []string{}); !ok {
				http.Error(w, "Invalid jwt token", http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}

	return h
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	// If the data implements the Public interface use it to prevent exposing sensitive data.
	if obj, ok := data.(Public); ok {
		data = obj
	}

	body, err := json.Marshal(data)
	must(err)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Decode(r *http.Request, v interface{}) error {
	// This can check if the OK method on a struct returns an error.
	// We can check if required fields are given this way.
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
