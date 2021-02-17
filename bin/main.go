package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	api "github.com/bastianhussi/todos-api"
	"github.com/bastianhussi/todos-api/profile"
)

var (
	ctx    = context.Background()
	config *api.Config
	router *mux.Router
	res    *api.Resources
	srv    *http.Server
)

func init() {
	var err error

	// if err would not been declared before config would only be shadow in this functions scope,
	// but then be nil in the main-functions scope
	config = api.NewConfig()

	res, err = api.NewResources(ctx, config)
	api.Must(err)

	router = mux.NewRouter().StrictSlash(true)

	// TODO: add profile route and use the auth adapter
	addHandle(router, []string{"/login"}, profile.NewLoginHandler(res.SharedKey), http.MethodPost)
	addHandle(router, []string{"/register"}, profile.NewRegisterHandler(), http.MethodPost)
	addHandle(router, []string{"/profile/{id}", "/p/{id}"}, api.Adapt(profile.NewProfileHandler(),
		api.Auth(res.SharedKey)), http.MethodGet,
		http.MethodPatch, http.MethodDelete)

	srv = &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.Port),
		WriteTimeout: config.Timeout.Write,
		ReadTimeout:  config.Timeout.Read,
		IdleTimeout:  config.Timeout.Idle,
		Handler:      api.Adapt(router, api.Recover(res.Logger), api.Logging(res.Logger)),
	}
}

// little helper functions that registers handles and adds the WithLogger and WithDB wrappers
//around them.
func addHandle(r *mux.Router, paths []string, h http.Handler, methods ...string) {
	for _, p := range paths {
		r.Handle(p, api.WithLogger(res.Logger, api.WithDB(res.DB, h))).Methods(methods...)
	}
}

func main() {
	defer res.Close()
	res.Logger.Printf("Server is running on %s 🚀\n", srv.Addr)

	// start the http server in a separate goroutine.
	go func() {
		if err := srv.ListenAndServe(); err == http.ErrServerClosed {
			res.Logger.Println("Server stopped 🛑")
		} else {
			res.Logger.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal to stop
	<-stop

	// create context with a timeout
	ctx, cancel := context.WithTimeout(ctx, config.Timeout.Shutdown)
	defer cancel()

	// try to shut the server down graceful by stop accepting incoming requests and finishing the
	//remaining ones. After the timeout finished kill the server.
	if err := srv.Shutdown(ctx); err != nil {
		res.Logger.Fatal(err)
	}
}
