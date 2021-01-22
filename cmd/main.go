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
)

var (
	l   *log.Logger
	srv *api.Server
)

func init() {
	logger := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)
	srv = api.NewServer(logger)
	res := api.NewResources(logger)
	srv.AddRoute(login.NewHandler(res))
	srv.AddRoute(register.NewHandler(res))
}

func main() {
	// c := new(api.Config)
	// if err := c.ReadConfig(); err != nil {
	// 	panic(err)
	// }
	// rdb := api.NewRedisClient()
	// defer rdb.Close()

	// db, err := api.NewDB()
	// if err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	go srv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
