package api

import (
	"net/http"
	"time"
)

// intialize an http server, generally app.Mount() will be passed as the mux
func (app *App) SetupHTTPServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         app.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
}

/*
create a mux server type & return to be run
all endpoints need to be defined in the mount function with their HTTP request
method/endpoint name and their corresponding HandleFunc
the root handler "/" must remain at the end of the function
*/
func (app *App) Mount() *http.ServeMux {
	mux := http.NewServeMux()
	app.Lg.Infof("new mux server created, setting up endpoint handlers")

	// define endpoints
	mux.HandleFunc("GET /about", app.HndlAbt)
	mux.HandleFunc("GET /health", app.HndlHealth)
	mux.HandleFunc("GET /bronto", app.HndlBronto)
	mux.HandleFunc("GET /bball", app.HndlBBall)
	mux.HandleFunc("GET /bball/about", app.HndlBBallAbt)
	mux.HandleFunc("GET /bball/seasons", app.HndlSeasons)
	mux.HandleFunc("GET /bball/teams", app.HndlTeams)
	mux.HandleFunc("GET /bball/player", app.HndlPlayer)
	mux.HandleFunc("GET /bball/games/recent", app.HndlRecentGames)
	mux.HandleFunc("GET /bball/league/scoring-leaders", app.HndlTopLgPlayers)
	mux.HandleFunc("GET /bball/teamrecs", app.HndlTeamRecords)

	// TESTING NEW ENDPOINTS 10/26/2025
	mux.HandleFunc("GET /bball/v2/players", app.HndlPlayerV2)

	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	mux.Handle("/docs/", http.HandlerFunc(app.ServeDocs))
	mux.Handle("/js/", http.HandlerFunc(app.JSNostore))
	mux.Handle("/css/", http.HandlerFunc(app.CSSNostore))
	mux.HandleFunc("/", app.HndlRoot)

	return mux
}
