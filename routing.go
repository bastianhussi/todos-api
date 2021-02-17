package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

type (
	key int

	logwrapper struct {
		logger  *log.Logger
		handler http.Handler
	}

	dbwrapper struct {
		dbSession *pg.DB
		handler   http.Handler
	}
)

const (
	dbKey     key = 0
	loggerKey key = 1
)

func WithLogger(l *log.Logger, h http.Handler) http.Handler {
	return &logwrapper{l, h}
}

func (l *logwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerKey, l.logger)

	l.handler.ServeHTTP(w, r.WithContext(ctx))
}

func WithDB(d *pg.DB, h http.Handler) http.Handler {
	return &dbwrapper{d, h}
}

// Provide a open db connection for each request using this and make sure the connection is closed
// when finished
func (d *dbwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := d.dbSession.Conn()
	defer conn.Close()

	ctx := r.Context()
	ctx = context.WithValue(ctx, dbKey, conn)

	d.handler.ServeHTTP(w, r.WithContext(ctx))
}

func LoggerFromContext(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(loggerKey).(*log.Logger)
	if !ok {
		panic("Could not receive the logger from the context of this request")
	}

	return logger
}

func DBFromContext(ctx context.Context) *pg.Conn {
	conn, ok := ctx.Value(dbKey).(*pg.Conn)
	if !ok {
		panic("Could not receive the database connection from the context of this request")
	}

	return conn
}

type Adapter func(http.Handler) http.Handler

func Logging(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			l.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
			// TODO: summarize response
			l.Printf("Request processed in %s\n", time.Since(start))
		})
	}
}

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

func Auth(sharedKey []byte) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			// FIXME: handle missing token
			token = strings.Split(token, "Bearer ")[1]

			if ok := VerifyJWT(token, sharedKey, []string{}); !ok {
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
