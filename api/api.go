package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

type Timing struct {
	CtxTimeout        time.Duration
	UpdateStoreTick   time.Duration
	UpdateStoreThresh time.Duration
	HealthCheckTick   time.Duration
	HealthCheckThreah time.Duration
}

type Env struct {
	PGEnv    *conn.DBEnv
	MongoEnv *conn.DBEnv
	SrvIP    string
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

func (e *Env) Load() error {
	pe, err := conn.Load("PG_HOST", "PG_PORT", "PG_USER", "PG_PASS", "PG_DB")
	if err != nil {
		return fmt.Errorf("failed to get postgres env: %v", err)
	}
	me, err := conn.Load("MONGO_HOST", "MONGO_PORT",
		"MONGO_INITDB_ROOT_USERNAME", "MONGO_INITDB_ROOT_PASSWORD", "MONGO_INITDB_DATABASE")
	if err != nil {
		return fmt.Errorf("failed to get mongodb env: %v", err)
	}
	ip := os.Getenv("SRV_IP")
	if ip == "" {
		return fmt.Errorf("error getting SRV_IP in env: %v", ip)
	}
	e.PGEnv = pe
	e.MongoEnv = me
	e.SrvIP = ip
	return nil
}

type Endpoints map[string]func(http.ResponseWriter, *http.Request)

// create a mux server type & return to be run
// all endpoints need to be defined in the mount function with their HTTP request
// method/endpoint name and their corresponding HandleFunc
// the root handler "/" must remain at the end of the function
func (app *App) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	app.ENDPOINTS = Endpoints{
		"GET /health":                       app.HndlHealth,
		"GET /dbhealth":                     app.HndlDBHealth,
		"GET /bball/seasons":                app.HndlSeasons,
		"GET /bball/teams":                  app.HndlTeams,
		"GET /bball/player":                 app.HndlPlayer,
		"GET /bball/games/recent":           app.HndlRecentGames,
		"GET /bball/league/scoring-leaders": app.HndlTopLgPlayers,
		"GET /bball/teamrecs":               app.HndlTeamRecords,
		"GET /bball/v2/players":             app.HndlPlayerV2,
	}

	for pattern, handler := range app.ENDPOINTS {
		mux.HandleFunc(pattern, handler)
	}

	app.Lg.Infof("%d endpoints mounted", len(app.ENDPOINTS))
	return mux
}

// intialize an http server, generally app.Mount() will be passed as the mux
func (app *App) SetupHTTPServer(mux *http.ServeMux) (*http.Server, error) {
	var ip string
	var ip_env string = "SRV_IP"
	ip = os.Getenv(ip_env)
	if ip == "" {
		return nil, fmt.Errorf("couldn't find %s in env", ip_env)
	}
	app.Addr = ip
	return &http.Server{
		Addr:         app.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}, nil
}
