package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg/v10"
)

// Router can add endpoints to the servers mux.
type Router interface {
	Route(mux *http.ServeMux)
}

// Resources are used by the Handlers to access databases, use logging utilities, ...
type Resources struct {
	Logger *log.Logger
	DB     *pg.DB
}

// NewResources creates a new Collection of resources.
func NewResources(l *log.Logger, db *pg.DB) *Resources {
	return &Resources{l, db}
}

// Logging monitors the incoming method and the time needed to process the request.
func (res *Resources) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res.Logger.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
		start := time.Now()
		defer next(w, r)
		res.Logger.Printf("Request processed in %s\n", time.Since(start))
	}
}

func (res *Resources) HandleRequest(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.Header().Add("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(code)
	_, err = w.Write([]byte(err.Error()))
	must(err)
}

func (res *Resources) HandleRequestPanic(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		res.Logger.Printf("Recovered from panic: %s\n", err)
		w.Header().Add("Content-Type", "plain/text; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte{})
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
