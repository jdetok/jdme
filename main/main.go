/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
)

func main() {
	app := &api.App{Started: false, QuickStart: true}

	app.SetupLoggers()
	cl, err := logd.SetupMongoLog()
	if err != nil {
		app.Lg.Fatalf("failed to connect to mongo: %v", err)
	}
	app.Lg.Mongo = cl
	defer func() {
		if err := app.Lg.Mongo.Disconnect(context.TODO()); err != nil {
			app.Lg.Fatalf("fatal mongo error: %v", err)
		}
	}()
	// persist file for quickstart
	app.SetupMemPersist("./persist_data/maps.json")

	if envErr := app.SetupEnv(); envErr != nil {
		app.Lg.Fatalf("fatal error reading .env file: %v", envErr)
	}

	// connect to bball postgres database
	if dbErr := app.SetupDB(); dbErr != nil {
		app.Lg.Fatalf("fatal db setup error: %v", dbErr)
	}

	// update Players, Seasons, Teams in memory structs
	memErrCh := make(chan error, 1)
	go func() {
		err := app.UpdateStore(app.QuickStart, 30*time.Minute)
		if err != nil {
			memErrCh <- err
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv := app.SetupHTTPServer(app.Mount())
	app.Lg.Infof("http server configured and endpoints mounted")

	// set the time for caching
	app.StartTime = time.Now()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.Lg.Errorf("http listen error occured: %v", err)
		}
	}()

	app.Lg.Infof("server listening at %v...\n", app.Addr)

	select {
	case <-shutdownCtx.Done():
		app.Lg.Quitf("shutdown signal received, shutting down...")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxTimeout); err != nil {
			app.Lg.Fatalf("Shutdown error: %v", err)
		}
	case err := <-memErrCh:
		app.Lg.Fatalf("memory refresh failed, exiting: %v", err)
	}
}
