package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/bastianhussi/todos-api"
	login "github.com/bastianhussi/todos-api/login"
	register "github.com/bastianhussi/todos-api/register"
	"github.com/go-pg/pg/v10"
)

var (
	l   *log.Logger
	c   *api.Config
	db  *pg.DB
	srv *api.Server
)

func init() {
	c, err := api.NewConfig()
	must(err)

	db, err = api.NewDB(c)
	must(err)

	l = log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)

	// TODO: use goroutines to handle these two tasks asyncronous
	res := api.NewResources(l, db)
	err = api.CreateSchema(db.Conn())
	must(err)

	srv = api.NewServer(l)
	srv.AddRoute(login.NewHandler(res))
	srv.AddRoute(register.NewHandler(res))

}

func main() {
	defer db.Close()
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