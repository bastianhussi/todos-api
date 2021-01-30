package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	logger *log.Logger
	router *mux.Router
	srv    *http.Server
}

func NewServer(l *log.Logger, c *Config) *Server {
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		WriteTimeout: c.Timeout.Write,
		ReadTimeout:  c.Timeout.Read,
		IdleTimeout:  c.Timeout.Idle,
		Handler:      router,
	}
	return &Server{l, router, srv}
}

// AddHandler takes in any struct capable of registering routes by adding new handles.
func (s *Server) AddHandler(patterns []string, h http.HandlerFunc, methods ...string) {
	for _, pattern := range patterns {
		s.router.HandleFunc(pattern, s.middleware(h)).Methods(methods...)
	}
}

// middleware is a wrapper function called before each incoming request.
// This function manages logging and recovering from panics.
func (s *Server) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle potential panics
		defer func() {
			if err := recover(); err != nil {
				s.logger.Println(err)
				http.Error(w, "Oh noo, something went wrong 🤯", http.StatusInternalServerError)
			}
		}()

		// Measure request time
		start := time.Now()
		s.logger.Printf("Got %s request at %s\n", r.Method, r.URL.Path)
		next(w, r)
		s.logger.Printf("Request processed in %s\n", time.Since(start))
	}
}

func (s *Server) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.Split(token, "Bearer ")[1]
		s.logger.Println(token)
		if ok := verifyToken(token, "secret", []string{}); !ok {
			http.Error(w, "Invalid jwt token", http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}

// Run starts the server. NOTE: This is a blocking function call.
// Therefore should be run in a separate goroutine.
func (s *Server) Run() {
	s.logger.Printf("Server is running on %s 🚀\n", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err == http.ErrServerClosed {
		s.logger.Println("Server stopped 🛑")
	} else {
		s.logger.Fatal(err)
	}
}

// Shutdown stops the running server. At first the server will have the possibility of shutting down
// gracefully. If the given context expires, the server is killed.
func (s *Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}
