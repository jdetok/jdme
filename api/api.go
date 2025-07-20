package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/cache"
)

type application struct {
	config     config
	database   *sql.DB
	StartTime  time.Time
	lastUpdate time.Time
	players    []cache.Player
	seasons    []cache.Season
	teams      []cache.Team
}

type config struct {
	addr string
	// cachePath string
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

	mux.HandleFunc("GET /about", app.abtHandler)
	mux.HandleFunc("GET /bronto", app.brontoHandler)
	mux.HandleFunc("GET /bball", app.bballHandler)
	mux.HandleFunc("GET /bball/about", app.bballAbtHandler)
	mux.HandleFunc("GET /bball/seasons", app.getSeasons)
	mux.HandleFunc("GET /bball/teams", app.getTeams)
	mux.HandleFunc("GET /bball/player", app.getPlayerDash)
	mux.HandleFunc("GET /bball/games/recent", app.getGamesRecentNew)

	mux.Handle("/js/", http.HandlerFunc(app.jsNoCache))
	mux.Handle("/css/", http.HandlerFunc(app.cssNoCache))
	mux.HandleFunc("/", app.rootHandler)

	return mux
}

func (app *application) setStartTime() {
	app.StartTime = time.Now()
}

func (app *application) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
