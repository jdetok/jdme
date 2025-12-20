/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/joho/godotenv"
)

func main() {
	app := &api.App{Started: false, QuickStart: false}
	errCh := make(chan error, 1)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("fatal error reading .env file: %v", err)
	}
	app.SetupLoggers()

	ml, err := logd.NewMongoLogger("log", "http")
	if err != nil {
		app.Lg.Fatalf("failed to connect to mongo: %v", err)
	}
	app.Lg.Mongo = ml
	// example players query: db.log.find({url: { $regex: "^/.*players.*$" } } )
	defer func() {
		if err := app.Lg.Mongo.Client.Disconnect(context.TODO()); err != nil {
			app.Lg.Fatalf("fatal mongo error: %v", err)
		}
	}()
	// persist file for quickstart
	app.SetupMemPersist("./persist_data/maps.json")

	// connect to bball postgres database
	if dbErr := app.SetupDB(); dbErr != nil {
		app.Lg.Fatalf("fatal db setup error: %v", dbErr)
	}

	// update Players, Seasons, Teams in memory structs

	go func() {
		err := app.UpdateStore(app.QuickStart, 30*time.Minute)
		if err != nil {
			errCh <- fmt.Errorf("in mem update error: %v", err)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv, err := app.SetupHTTPServer(app.Mount())
	if err != nil {
		app.Lg.Fatalf("failed to setup HTTP server: %v", err)
	}

	// set the time for caching
	app.StartTime = time.Now()

	go func() {
		app.Lg.Infof("starting server at %v...\n", app.Addr)
		err := srv.ListenAndServe()
		// if err != nil && err != http.ErrServerClosed {
		if err != nil {
			errCh <- fmt.Errorf("http listen error occured: %v", err)
		}
		app.Lg.Infof("server started at %s\n", app.Addr)
	}()

	select {
	case <-shutdownCtx.Done():
		app.Lg.Quitf("shutdown signal received, shutting down...")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxTimeout); err != nil {
			app.Lg.Fatalf("Shutdown error: %v", err)
		}
	case err := <-errCh:
		app.Lg.Fatalf("fatal error occurred: %v", err)
	}
}
