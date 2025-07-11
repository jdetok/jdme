package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jdetok/go-api-jdeko.me/internal/store"
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
	addr      string
	cachePath string
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

	// ENDPOINTS 06/19
	mux.HandleFunc("GET /devl/bball", app.bballDevHandler)
	mux.HandleFunc("GET /devl/bball/player", app.getPlayerDash)
	mux.HandleFunc("GET /bball/games/recent", app.getGamesRecentNew)
	mux.HandleFunc("GET /bball/games/recent/top-scorer", app.getTopScorerNew)
	mux.HandleFunc("GET /bball/players", app.getStats)
	mux.HandleFunc("GET /bball/players/id", app.getPlayerId)
	mux.HandleFunc("GET /bball/players/random", app.getRandPlayer)
	mux.HandleFunc("GET /bball/seasons", app.getSeasons)
	mux.HandleFunc("GET /bball/teams", app.getTeams)
	mux.HandleFunc("GET /bball", app.bballHandler)

	// SERVES STATIC SITE IN WEB DIRECTORY, DON'T CACHE JS & CSS
	// mux.Handle("/devl/", http.HandlerFunc(app.devHandler))

	// mux.HandleFunc("/dev", app.devHandler)
	mux.Handle("/devl/", http.HandlerFunc(app.devHandler))
	mux.Handle("/js/", http.HandlerFunc(app.jsNoCache))
	mux.Handle("/css/", http.HandlerFunc(app.cssNoCache))
	// mux.Handle("/devl/css/", http.HandlerFunc(app.cssDevNoCache))
	mux.HandleFunc("/", app.rootHandler)

	// call app.mount() to get mux then app.run(mux) to run server
	return mux
}

func (app *application) setStartTime() {
	app.StartTime = time.Now()
}

func (app *application) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
