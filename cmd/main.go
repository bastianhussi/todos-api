package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"

	api "github.com/bastianhussi/todos-api"
	"github.com/bastianhussi/todos-api/login"
	"github.com/bastianhussi/todos-api/register"
)

var (
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

	loginHandle := api.WithLogging(logger, api.WithDB(db, login.NewHandler()))
	registerHandle := api.WithLogging(logger, api.WithDB(db, register.NewHandler()))

	router.Handle("/login", loginHandle).Methods(http.MethodPost)
	router.Handle("/register", registerHandle).Methods(http.MethodPost)

	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		WriteTimeout: config.Timeout.Write,
		ReadTimeout:  config.Timeout.Read,
		IdleTimeout:  config.Timeout.Idle,
		Handler:      api.Adapt(router, api.Recover(logger), api.Logging(logger)),
	}
}

func main() {
	defer db.Close()

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

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}
