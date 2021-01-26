package api

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	Resources struct {
		Logger *log.Logger
		DB     *pg.DB
	}
	Server struct {
		Res    *Resources
		Router *mux.Router
		srv    *http.Server
	}
)

func NewResources(c *Config)  (*Resources, error) {
	db, err := NewDB(c)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)

	return &Resources{
		logger, db,
	}, nil
}

func NewServer(c *Config, r *Resources) *Server {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", c.Port),
		ReadTimeout:  c.Timeout.Read,
		WriteTimeout: c.Timeout.Write,
		IdleTimeout:  c.Timeout.Shutdown,
		Handler:      router,
	}

	return &Server{
		r,
		router,
		server,
	}
}

func (s *Server) AddRoute(paths []string, handler http.HandlerFunc, methods ...string) {
	for _, path := range paths {
		s.Router.HandleFunc(path, s.middleware(handler)).Methods(methods...)
	}
}

// Logging monitors the incoming method and the time needed to process the request.
func (s *Server) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle potential panics
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Oh noo, something went wrong 🤯", http.StatusInternalServerError)
			}
		}()

		log := s.Res.Logger

		// Measure request time
		start := time.Now()
		log.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Request processed in %s\n", time.Since(start))
	}
}

func (s *Server) Run() {
	log := s.Res.Logger;

	log.Printf("Server is running on %s 🚀\n", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err == http.ErrServerClosed {
		log.Println("Server stopped 🛑")
	} else {
		log.Fatal(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.Res.Logger.Fatal(err)
	}
}
