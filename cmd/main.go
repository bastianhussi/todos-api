package main

import (
	"context"
	"github.com/bastianhussi/todos-api/login"
	"github.com/bastianhussi/todos-api/profile"
	"github.com/bastianhussi/todos-api/register"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/bastianhussi/todos-api"
)

var (
	c   *api.Config
	res *api.Resources
	s   *api.Server
)

func init() {
	c, err := api.NewConfig()
	must(err)

	// TODO: use goroutines to handle these two tasks asyncronous
	res, err = api.NewResources(c)
	must(err)
	err = api.CreateSchema(res.DB.Conn())
	must(err)

	s = api.NewServer(c, res)
	// Doesn't work because Mux.HandleFunc is not a function, its a method of *mux.Router
	login.NewHandler(s)
	register.NewHandler(s)
	profile.NewHandler(s)
}

func main() {
	defer res.DB.Close()
	go s.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
