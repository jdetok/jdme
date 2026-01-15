package api

import (
	"os"
	"time"

	"github.com/jdetok/jdme/pkg/logd"
	"github.com/jdetok/jdme/pkg/memd"
	"github.com/jdetok/jdme/pkg/pgdb"
)

type Timing struct {
	CtxTimeout        time.Duration
	UpdateStoreTick   time.Duration
	UpdateStoreThresh time.Duration
	HealthCheckTick   time.Duration
	HealthCheckThreah time.Duration
}

// GLOBAL APP STRUCT
type App struct {
	E          Env
	T          Timing
	ENDPOINTS  Endpoints
	Addr       string
	DB         pgdb.DB
	DBConf     pgdb.DBConfig
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
