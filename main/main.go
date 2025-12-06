/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/golib/envd"
)

func main() {
	app := &api.App{}
	app.Started = false
	var quickstart bool = false

	// main logger
	f, err := logd.SetupLogdF("./z_log/applog")
	if err != nil {
		log.Fatal(err)
	}

	// debug logger
	df, err := logd.SetupLogdF("./z_log/debug_applog")
	if err != nil {
		log.Fatal(err)
	}

	// http logger
	hf, err := logd.SetupLogdF("./z_log/http_log")
	if err != nil {
		log.Fatal(err)
	}

	// Logd setup with each logger
	app.Lg = logd.NewLogd(io.MultiWriter(os.Stdout, f), df, hf)
	app.Lg.Infof("started app and created log file")

	// persist file for quickstart
	fp := "./persist_data/maps.json"
	persistP, err := filepath.Abs(fp)
	if err != nil {
		app.Lg.Fatalf("failed to get absolute path of %s\n**%v\n", fp, err)
	}
	app.MStore.PersistPath = persistP
	app.Lg.Warnf("path: %s", app.MStore.PersistPath)

	err = envd.LoadDotEnv()
	if err != nil {
		app.Lg.Fatalf("failed to load variables in .env file to env\n%v", err)
	}
	hostaddr, err := envd.GetEnvStr("SRV_IP")
	if err != nil {
		app.Lg.Fatalf("failed to get server IP from .env\n%v", err)
	}
	app.Addr = hostaddr

	// connect to bball postgres database
	db, err := pgdb.PostgresConn()
	if err != nil {
		app.Lg.Fatalf("failed to create connection to postgres\n%v", err)
	}
	app.DB = db

	// update Players, Seasons, Teams in memory structs
	go func(*api.App) {
		err := app.UpdateStore(quickstart, 120*time.Minute)
		if err != nil {
			app.Lg.Fatalf("error updating store: %v", err)
		}
	}(app)

	// mount mux server, sets up all endpoint handlers
	mux := app.Mount()
	app.Lg.Infof("http mux server mounted, starting server")

	if err := app.RunGraceful(mux); err != nil {
		app.Lg.Fatalf("FATAL server failed to run\n%v", err)
	}

}
