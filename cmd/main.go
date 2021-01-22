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
)

var (
	l *log.Logger
	s *api.Server
)

func init() {
	l := log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)
	s = api.NewServer(l)
	s.AddRoute(login.NewHandler(l))
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

	go s.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
