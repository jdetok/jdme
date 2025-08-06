package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/store"
)

type application struct {
	config     config
	database   *sql.DB
	StartTime  time.Time
	lastUpdate time.Time
	players    []store.Player
	seasons    []store.Season
	teams      []store.Team
}

type config struct {
	addr string
	// storePath string
}

func (app *application) run(mux *http.ServeMux) error {

	// server configuration
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// set the time for caching
	app.setStartTime()
	fmt.Printf("http server configured and starting at %v...\n",
		app.StartTime.Format("2006-01-02 15:04:05"))

	return srv.ListenAndServe()
}

// returns type ServeMux for a router
func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	// standardize handlers: end with Hndl e.g. abtHndl, brontoHndl
	mux.HandleFunc("GET /about", app.abtHndl)
	mux.HandleFunc("GET /bronto", app.brontoHndl)
	mux.HandleFunc("GET /bball", app.bballHndl)
	mux.HandleFunc("GET /bball/about", app.bballAbtHndl)
	mux.HandleFunc("GET /bball/seasons", app.seasonsHndl)
	mux.HandleFunc("GET /bball/teams", app.teamsHndl)
	mux.HandleFunc("GET /bball/player", app.playerDashHndl)
	mux.HandleFunc("GET /bball/games/recent", app.recGameHndl)

	mux.Handle("/js/", http.HandlerFunc(app.jsNostore))
	mux.Handle("/css/", http.HandlerFunc(app.cssNostore))
	mux.HandleFunc("/", app.rootHndl)

	return mux
}

func (app *application) setStartTime() {
	app.StartTime = time.Now()
}

func (app *application) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
