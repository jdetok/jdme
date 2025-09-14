package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

/*
main struct referenced through the app. contains configs, , pool,
in-memory player, season, team slices
*/
type App struct {
	Config     Config
	Database   *sql.DB
	StartTime  time.Time
	LastUpdate time.Time
	Started    uint8
	Players    []Player
	Seasons    []Season
	Teams      []Team
}

// configs, currently only contains server address
type Config struct {
	Addr string
}

// accept slice of bytes in JSON structure and write to response writers
func (app *App) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
create a mux server type & return to be run
all endpoints need to be defined in the mount function with their HTTP request
method/endpoint name and their corresponding HandleFunc
the root handler "/" must remain at the end of the function
*/
func (app *App) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	// define endpoints
	mux.HandleFunc("GET /about", app.AboutHndl)
	mux.HandleFunc("GET /bronto", app.BrontoHndl)
	mux.HandleFunc("GET /bball", app.BBallHndl)
	mux.HandleFunc("GET /bball/about", app.BballAbtHndl)
	mux.HandleFunc("GET /bball/seasons", app.SeasonsHndl)
	mux.HandleFunc("GET /bball/teams", app.TeamsHndl)
	mux.HandleFunc("GET /bball/player", app.PlayerDashHndl)
	mux.HandleFunc("GET /bball/games/recent", app.RecentGameHndl)

	// serve static files
	mux.Handle("/js/", http.HandlerFunc(app.JSNostore))
	mux.Handle("/css/", http.HandlerFunc(app.CSSNostore))
	mux.HandleFunc("/", app.RootHndl)

	return mux
}

// runs the http server - must be called after mount is successfully executed
func (app *App) Run(mux *http.ServeMux) error {

	// server configuration
	srv := &http.Server{
		Addr:         app.Config.Addr,
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
