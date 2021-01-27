package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bastianhussi/todos-api/login"
	"github.com/bastianhussi/todos-api/register"
	"github.com/gorilla/mux"

	api "github.com/bastianhussi/todos-api"
)

var (
	res    *api.Resources
	router *mux.Router
	srv    *api.Server
)

func init() {
	config, err := api.NewConfig()
	must(err)

	// TODO: use goroutines to handle these two tasks asyncronous
	res, err = api.NewResources(config)
	must(err)
	must(api.CreateSchema(res.DB.Conn()))
	srv = api.NewServer(res.Logger, config)
	login.NewHandler(res).RegisterRoute(srv)
	register.NewHandler(res).RegisterRoute(srv)
}

func main() {
	defer res.DB.Close()
	go srv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
