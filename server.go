package api

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Mux    *http.ServeMux
	srv    *http.Server
	logger *log.Logger
}

func NewServer(l *log.Logger) *Server {
	s := new(Server)
	s.logger = l
	s.Mux = http.NewServeMux()
	s.srv = &http.Server{
		Addr:        ":3000",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler: s.Mux,
	}

	return s
}

func (s *Server) AddRoute(r Router) {
	r.Route(s.Mux)
}

func (s *Server) Run() {
	if err := s.srv.ListenAndServe(); err == http.ErrServerClosed {
		s.logger.Println("Server stopped 🛑")
	} else {
		s.logger.Fatal(err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}
