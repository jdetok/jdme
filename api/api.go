package api

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
)

// GLOBAL APP STRUCT
type App struct {
	Addr       string
	DB         *sql.DB
	StartTime  time.Time
	LastUpdate time.Time
	Started    bool
	Store      memd.InMemStore
	MStore     memd.MapStore
	Logf       *os.File
	QLogf      *os.File
	Lg         *logd.Logd
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
	// mux.HandleFunc("/docs/", app.HndlWiki)
	// serve static files
	mux.Handle("/js/", http.HandlerFunc(app.JSNostore))
	mux.Handle("/css/", http.HandlerFunc(app.CSSNostore))
	mux.HandleFunc("/", app.HndlRoot)

	return mux
}

// runs the http server - must be called after mount is successfully executed
func (app *App) Run(mux *http.ServeMux) error {

	// server configuration
	srv := &http.Server{
		Addr:         app.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// set the time for caching
	app.StartTime = time.Now()

	app.Lg.Infof("http server configured | running server at %v...\n", app.Addr)

	// run the HTTP server
	return srv.ListenAndServe()
}

func (app *App) RunGraceful(mux *http.ServeMux) error {
	// server configuration
	srv := &http.Server{
		Addr:         app.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// set the time for caching
	app.StartTime = time.Now()
	app.Lg.Infof("http server configured | running server at %v...\n", app.Addr)

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			app.Lg.Errorf("a listen error occured")
		} // send regardless so channel doesn't block
		errCh <- err
	}()

	// listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		app.Lg.Infof("shutdown signal received: %v", sig)
		return nil
	case err := <-errCh:
		if err != nil {
			app.Lg.Errorf("server stopped with error: %v", err)
			return err
		}
	}
	return nil
	// run the HTTP server
	// return srv.ListenAndServe()
}
