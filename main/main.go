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
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/golib/envd"
)

func main() {

	app := &api.App{}

	// logger setup - opens a *os.File which implements io writer interface
	f, err := logd.SetupLogdF("./z_log/applog")
	if err != nil {
		log.Fatal(err)
	}
	app.Logf = f

	// SETUP MAIN APP LOGGER
	app.Lg = logd.NewLogd(io.MultiWriter(os.Stdout, f))

	// log file created confirmation
	app.Lg.Infof("started app and created log file")

	// load environment variables from .env file
	err = envd.LoadDotEnv()
	if err != nil {
		app.Lg.Fatalf("failed to load variables in .env file to env\n%v", err)
	}

	// set the server IP address
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

	// build new map store
	app.MStore.Maps = &memd.StMaps{}
	app.MStore.Setup(app.DB)
	if err := app.MStore.Rebuild(app.DB, app.Lg); err != nil {
		app.Lg.Fatalf("failed to build in memory map stores")
	}
	app.Lg.Infof("in memory map store setup complete")

	// set started = 0 so first check to update store runs setups
	app.Started = false

	// update Players, Seasons, Teams in memory structs
	go app.CheckInMemStructs(30*time.Second, 300*time.Second)

	// mount mux server, sets up all endpoint handlers
	mux := app.Mount()
	app.Lg.Infof("http mux server mounted, starting server")

	// run the http mux server
	if err := app.Run(mux); err != nil {
		app.Lg.Fatalf("FATAL server failed to run\n%v", err)
	}
}
