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
	"sync"
	"syscall"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/mgo"
	"github.com/joho/godotenv"
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
	app.MStore.PersistPath = "./persist/maps.json"

	if err := godotenv.Load(); err != nil {
		fmt.Printf("fatal error reading .env file: %v\n", err)
		os.Exit(1)
	}

	app.SetupLoggers()
	ml, err := mgo.NewMongoLogger("log", "http")
	if err != nil {
		app.Lg.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		if err := ml.Client.Disconnect(context.TODO()); err != nil {
			app.Lg.Fatalf("fatal mongo error: %v", err)
		}
	}()

	app.Lg.Infof("environment variables loaded and loggers setup successfully")

	app.Lg.Mongo = ml
	if dbErr := app.SetupDB(); dbErr != nil {
		app.Lg.Fatalf("fatal db setup error: %v", dbErr)
	}

	app.Lg.Infof("database connection created successfully")

	srv, err := app.SetupHTTPServer(app.Mount())
	if err != nil {
		app.Lg.Fatalf("failed to setup HTTP server: %v", err)
	}

	app.Lg.Infof("HTTP server configured successfully")

	shutdownCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 3)
	var wg sync.WaitGroup
	wg.Go(func() {
		defer wg.Done()
		app.Lg.Infof("starting HTTP server at %v...\n", app.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if shutdownCtx.Err() == nil {
				select {
				case errCh <- fmt.Errorf("http listen error occured: %v", err):
				default:
					app.Lg.Infof("HTTP server listening at %v...\n", app.Addr)
				}
			}
		}
	})
	wg.Go(func() {
		defer wg.Done()
		url := fmt.Sprintf("http://%s/health", srv.Addr)
		var lastCheck time.Time
		thresh := 60 * time.Second
		ticker := time.NewTicker(time.Second)
		var fails []struct{}
		failLimit := 5
		defer ticker.Stop()
		for range ticker.C {
			if time.Since(lastCheck) >= thresh {
				lastCheck = time.Now()
				resp, err := http.Get(url)
				if err != nil {
					fails = append(fails, struct{}{})
					app.Lg.Errorf("health check failed: %v", err)
					if len(fails) < failLimit {
						continue
					}
					if shutdownCtx.Err() == nil {
						select {
						case errCh <- fmt.Errorf("error occured during health check: %v", err):
						default:
						}
						return
					}
				}
				defer resp.Body.Close()
				switch resp.StatusCode {
				case http.StatusOK:
					app.Lg.Infof("health check passed: %s", resp.Status)
				default:
					app.Lg.Infof("health check status: %s", resp.Status)

				}
			}
		}
	})

	wg.Go(func() { // update in memory store every
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				select {
				case errCh <- fmt.Errorf("UpdateStore panic: %v", r):
				default:
				}
			}
		}()
		err := app.UpdateStore(shutdownCtx, app.QuickStart,
			app.T.UpdateStoreTick, app.T.UpdateStoreThresh)
		if err != nil {
			if shutdownCtx.Err() == nil {
				select {
				case errCh <- fmt.Errorf("in mem update error: %w", err):
				default:
				}
			}
		}
	})

	select {
	case <-shutdownCtx.Done():
		app.Lg.Quitf("shutdown signal received, shutting down...")

	case err := <-errCh:
		app.Lg.Errorf("fatal error occurred: %v", err)
		stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), app.T.CtxTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Lg.Errorf("shutdown error: %v", err)
	}

	wg.Wait()
	app.Lg.Infof("shutdown complete")
}
