/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	ENV_FILE        = ".env"
	PERSIST_FILE    = "./persist/maps.json"
	APPLOG_FILE     = "./z_log/app/applog"
	DEBUG_FILE      = "./z_log/dbg/debug"
	MONGO_LOG_DB    = "log"
	MONGO_HTTP_COLL = "http"
)

func main() {
	app := &api.App{Started: false, QuickStart: true}
	app.T = api.Timing{
		CtxTimeout:        10 * time.Second,
		UpdateStoreTick:   1 * time.Minute,
		UpdateStoreThresh: 30 * time.Minute,
		HealthCheckTick:   1 * time.Second,
		HealthCheckThreah: 120 * time.Second,
	}
	app.MStore.PersistPath = PERSIST_FILE

	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("fatal error reading .env file: %v\n", err)
		os.Exit(1)
	}

	l, err := logd.SetupLoggers(APPLOG_FILE, DEBUG_FILE, MONGO_LOG_DB, MONGO_HTTP_COLL)
	if err != nil {
		fmt.Printf("fatal error setting up loggers: %v\n", err)
		os.Exit(1)
	}
	app.Lg = l
	defer func() { // disconnect from mongo when app shuts down
		if err := l.Mongo.Client.Disconnect(context.TODO()); err != nil {
			app.Lg.Fatalf("fatal mongo error: %v", err)
		}
	}()
	app.Lg.Infof("environment variables loaded from file: %s | loggers setup successfully", ENV_FILE)

	db, err := pgdb.PostgresConn()
	if err != nil {
		app.Lg.Errorf("failed to create connection to postgres\n%v", err)
	}
	app.DB = db

	app.Lg.Infof("database connection created successfully")

	srv, err := app.SetupHTTPServer(app.Mount())
	if err != nil {
		app.Lg.Fatalf("failed to setup HTTP server: %v", err)
	}

	app.Lg.Infof("HTTP server configured successfully")

	sigCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// errCh := make(chan error, 3)
	// var wg sync.WaitGroup
	g, ctx := errgroup.WithContext(sigCtx)
	g.Go(func() error {
		// defer wg.Done()
		app.Lg.Infof("starting HTTP server at %v...\n", app.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	g.Go(func() error {
		// defer wg.Done()
		url := fmt.Sprintf("http://%s/health", srv.Addr)
		var lastCheck time.Time

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var fails int
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if time.Since(lastCheck) >= app.T.HealthCheckThreah {
					lastCheck = time.Now()
					resp, err := http.Get(url)
					if err != nil {
						fails++
						app.Lg.Errorf("health check failed: %v", err)
						if fails <= 5 {
							continue
						}
						return err
					}
					resp.Body.Close()
				}
			}
		}
	})

	g.Go(func() error { // update in memory store every
		// defer wg.Done()
		// defer func() error {
		// 	if r := recover(); r != nil {
		// 		return fmt.Errorf("UpdateStore panic: %v", r)
		// 	}
		// }()
		err := app.UpdateStore(ctx, app.QuickStart,
			app.T.UpdateStoreTick, app.T.UpdateStoreThresh)
		if err != nil {
			if ctx.Err() == nil {
				return fmt.Errorf("in mem update error: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		app.Lg.Infof("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		app.Lg.Errorf("wait error: %v", err)
	}
	app.Lg.Infof("shutdown complete")
}
