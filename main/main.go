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
	PROD_URL        = "https://jdeko.me/"
	ENV_FILE        = ".env"
	PERSIST_FILE    = "./persist/maps.json"
	APPLOG_FILE     = "./z_log/app/applog"
	DEBUG_FILE      = "./z_log/dbg/debug"
	MONGO_LOG_DB    = "log"
	MONGO_HTTP_COLL = "http"
	PG_OPEN         = 80
	PG_IDLE         = 30
	PG_LIFE         = 30
	QUICKSTART      = false
	IS_PROD         = false
)

func main() {
	app := &api.App{Started: false, QuickStart: QUICKSTART}
	app.T = api.Timing{
		CtxTimeout:        10 * time.Second,
		UpdateStoreTick:   1 * time.Minute,
		UpdateStoreThresh: 2 * time.Minute,
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

	// database config
	app.DBConf = *pgdb.NewDBConf(PG_OPEN, PG_IDLE, PG_LIFE*time.Minute)
	db, err := pgdb.PostgresConn(&app.DBConf)
	if err != nil {
		app.Lg.Errorf("failed to create connection to postgres\n%v", err)
	}
	app.DB = db
	app.Lg.Infof("database connection created successfully")

	// http server config
	srv, err := app.SetupHTTPServer(app.Mount())
	if err != nil {
		app.Lg.Fatalf("failed to setup HTTP server: %v", err)
	}

	app.Lg.Infof("HTTP server configured successfully")

	sigCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	g, ctx := errgroup.WithContext(sigCtx)
	g.Go(func() error { // LISTEN FOR HTTP REQUESTS
		app.Lg.Infof("starting HTTP server at %v...\n", app.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	g.Go(func() error { // RUN HEALTH CHECKS
		var lastCheck time.Time
		healthChk := "health"
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
					var url string
					// if IS_PROD {
					// 	url = PROD_URL + healthChk
					// } else {
					// 	url = fmt.Sprintf("http://%s/%s", srv.Addr, healthChk)
					// }
					url = fmt.Sprintf("http://%s/%s", srv.Addr, healthChk)

					resp, err := http.Get(url)
					if err != nil {
						fails++
						app.Lg.Errorf("health check to %s failed: %v", url, err)
						if fails <= 5 {
							continue
						}
						return err
					}
					app.Lg.Infof("received passing healthcheck from %s | %v", url, resp.Status)
					resp.Body.Close()
				}
			}
		}
	})

	g.Go(func() error { // update in memory store every
		err := app.UpdateStore(ctx, app.QuickStart,
			app.T.UpdateStoreTick, app.T.UpdateStoreThresh)
		if err != nil {
			if ctx.Err() == nil {
				return fmt.Errorf("in mem update error: %w", err)
			}
			return err
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
