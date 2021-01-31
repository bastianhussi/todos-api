package main

import (
	"context"
	"fmt"
	"github.com/bastianhussi/todos-api/register"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "github.com/bastianhussi/todos-api"
	"github.com/bastianhussi/todos-api/login"
	"github.com/bastianhussi/todos-api/profile"
)

var (
	config *api.Config
	router *mux.Router
	srv    *http.Server
	logger *log.Logger
	db     *pg.DB
)

var ctx = context.Background()

func init() {
	config, err := api.NewConfig()
	api.Must(err)

	db, err = api.NewDB(ctx, config)
	api.Must(err)
	conn := db.Conn()
	defer conn.Close()

	logger := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)

	api.Must(api.CreateSchema(conn))

	router := mux.NewRouter().StrictSlash(true)

	// TODO: add profile route and use the auth adapter
	addHandle(router, []string{"/login"}, login.NewHandler(config.SharedKey), http.MethodPost)
	addHandle(router, []string{"/register"}, register.NewHandler(), http.MethodPost)
	addHandle(router, []string{"/profile/{id}", "/p/{id}"}, api.Adapt(profile.NewHandler(),
		api.Auth()), http.MethodGet,
		http.MethodPatch, http.MethodDelete)

	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		WriteTimeout: config.Timeout.Write,
		ReadTimeout:  config.Timeout.Read,
		IdleTimeout:  config.Timeout.Idle,
		Handler:      api.Adapt(router, api.Recover(logger), api.Logging(logger)),
	}
}

// little helper functions that registers handles and adds the WithLogger and WithDB wrappers
//around them.
func addHandle(r *mux.Router, paths []string, h http.Handler, methods ...string) {
	for _, p := range paths {
		r.Handle(p, api.WithLogger(logger, api.WithDB(db, h))).Methods(methods...)
	}
}

func main() {
	defer db.Close()

	// start the http server in a separate goroutine.
	go func() {
		logger.Printf("Server is running on %s 🚀\n", srv.Addr)
		if err := srv.ListenAndServe(); err == http.ErrServerClosed {
			logger.Println("Server stopped 🛑")
		} else {
			logger.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal to stop
	<-stop

	// create context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout.Shutdown)
	defer cancel()

	// try to shut the server down graceful by stop accepting incoming requests and finishing the
	//remaining ones. After the timeout finished kill the server.
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}