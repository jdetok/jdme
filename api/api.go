package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/golib/envd"
)

// GLOBAL APP STRUCT
type App struct {
	ENDPOINTS  Endpoints
	Addr       string
	DB         pgdb.DB
	StartTime  time.Time
	LastUpdate time.Time
	Started    bool
	QuickStart bool
	Store      memd.InMemStore
	MStore     memd.MapStore
	Logf       *os.File
	QLogf      *os.File
	Lg         *logd.Logd
}

func (app *App) SetupLoggers() {
	// main logger
	f, err := logd.SetupLogdF("./z_log/app/applog")
	if err != nil {
		log.Fatal(err)
	}

	// debug logger
	df, err := logd.SetupLogdF("./z_log/dbg/debug")
	if err != nil {
		log.Fatal(err)
	}

	// Logd setup with each logger
	app.Lg = logd.NewLogd(io.MultiWriter(os.Stdout, f), df)
}

func (app *App) SetupMemPersist(fp string) {
	persistP, err := filepath.Abs(fp)
	if err != nil {
		app.Lg.Fatalf("failed to get absolute path of %s\n**%v\n", fp, err)
	}
	app.MStore.PersistPath = persistP
	app.Lg.Infof("mem persist path: %s", app.MStore.PersistPath)
}

func (app *App) SetupEnv() error {
	err := envd.LoadDotEnv()
	if err != nil {
		return fmt.Errorf("failed to load variables in .env file to env\n%v", err)
	}
	hostaddr, err := envd.GetEnvStr("SRV_IP")
	if err != nil {
		return fmt.Errorf("failed to get server IP from .env\n%v", err)
	}
	app.Addr = hostaddr
	return nil
}

func (app *App) SetupDB() error {
	db, err := pgdb.PostgresConn()
	if err != nil {
		return fmt.Errorf("failed to create connection to postgres\n%v", err)
	}
	app.DB = db
	return nil
}
