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

	// logger setup - opens a *os.File which implements io writer interface
	f, err := logd.SetupLogdF("./z_log/applog")
	if err != nil {
		log.Fatal(err)
	}
	app.Logf = f
	df, err := logd.SetupLogdF("./z_log/debug_applog")
	if err != nil {
		log.Fatal(err)
	}
	app.QLogf = df

	hf, err := logd.SetupLogdF("./z_log/http_log")
	if err != nil {
		log.Fatal(err)
	}
	// SETUP MAIN APP LOGGER
	app.Lg = logd.NewLogd(io.MultiWriter(os.Stdout, f), df, hf)

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

	fp := "./persist_data/maps.json"
	persistP, err := filepath.Abs(fp)
	if err != nil {
		app.Lg.Fatalf("failed to get absolute path of %s\n**%v\n", fp, err)
	}
	app.MStore.PersistPath = persistP

	// if err := app.MStore.SetupFromPersist(); err != nil {
	// 	app.Lg.Fatalf("failed to build in memory map stores: %v", err)
	// }
	// app.Lg.Infof("in memory map store setup complete")
	// app.MStore.Persist()

	// set started = 0 so first check to update store runs setups
	app.Started = false

	// update Players, Seasons, Teams in memory structs
	// go app.CheckInMemStructs(300*time.Second, 30*time.Second)
	go func(*api.App) {
		err := app.UpdateStore(false, 300*time.Second, 30*time.Second)
		if err != nil {
			app.Lg.Fatalf("error updating store: %v", err)
		}
	}(app)
	// go app.CheckInMemStructs(300*time.Second, 30*time.Second)

	// go func(*sync.WaitGroup, *api.App) {
	// 	if err := app.MStore.Rebuild(app.DB, app.Lg); err != nil {
	// 		app.Lg.Errorf("failed to update player map")
	// 	}
	// }(wg, app)

	// mount mux server, sets up all endpoint handlers
	mux := app.Mount()
	app.Lg.Infof("http mux server mounted, starting server")

	if err := app.RunGraceful(mux); err != nil {
		app.Lg.Fatalf("FATAL server failed to run\n%v", err)
	}

}
