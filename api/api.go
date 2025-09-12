package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/store"
)

/*
main struct referenced through the app. contains configs, database pool,
in-memory player, season, team slices
*/
type application struct {
	config     config
	database   *sql.DB
	StartTime  time.Time
	lastUpdate time.Time
	players    []store.Player
	seasons    []store.Season
	teams      []store.Team
}

// configs, currently only contains server address
type config struct {
	addr string
}

func (app *application) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
create a mux server type & return to be run
all endpoints need to be defined in the mount function with their HTTP request
method/endpoint name and their corresponding HandleFunc
the root handler "/" must remain at the end of the function
*/
func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	// define endpoints
	mux.HandleFunc("GET /about", app.abtHndl)
	mux.HandleFunc("GET /bronto", app.brontoHndl)
	mux.HandleFunc("GET /bball", app.bballHndl)
	mux.HandleFunc("GET /bball/about", app.bballAbtHndl)
	mux.HandleFunc("GET /bball/seasons", app.seasonsHndl)
	mux.HandleFunc("GET /bball/teams", app.teamsHndl)
	mux.HandleFunc("GET /bball/player", app.playerDashHndl)
	mux.HandleFunc("GET /bball/games/recent", app.recGameHndl)

	// serve static files
	mux.Handle("/js/", http.HandlerFunc(app.jsNostore))
	mux.Handle("/css/", http.HandlerFunc(app.cssNostore))
	mux.HandleFunc("/", app.rootHndl)

	return mux
}

// runs the http server - must be called after mount is successfully executed
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
	app.StartTime = time.Now()

	fmt.Printf("http server configured and starting at %v...\n",
		app.StartTime.Format("2006-01-02 15:04:05"))

	// run the HTTP server
	return srv.ListenAndServe()
}
