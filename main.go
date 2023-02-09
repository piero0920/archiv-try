package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/piero0920/archiv-try/docs"
	"github.com/piero0920/archiv-try/pkg/chatlogger"
	"github.com/piero0920/archiv-try/pkg/cronjobs"
	"github.com/piero0920/archiv-try/pkg/database"
	"github.com/piero0920/archiv-try/pkg/helpers"
	"github.com/piero0920/archiv-try/pkg/logger"
	"github.com/piero0920/archiv-try/pkg/router"
)

func main() {
	// wait for db
	if db, err := database.DB.DB(); err != nil {
		for {
			if e := db.Ping(); e == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// start http server
	srv := &http.Server{
		Addr:    "0.0.0.0:5000",
		Handler: router.Init(),
	}

	// import settings from os env to database
	if err := helpers.ImportEnvToDb(); err != nil {
		panic(err)
	}

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Fatalf("[main] listen: %s\n", err)
		}
	}()

	// start cronjobs
	if err := cronjobs.Init(); err != nil {
		panic(err)
	}

	// run chatlogger
	go chatlogger.Run()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug.Println("[main] Shutting down server...")
	database.Close()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error.Fatal("[main] Server forced to shutdown: ", err)
	}

	logger.Debug.Println("[main] Server exiting")
}
