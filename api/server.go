package api

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Endpoints map[string]func(http.ResponseWriter, *http.Request)

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

// create a mux server type & return to be run
// all endpoints need to be defined in the mount function with their HTTP request
// method/endpoint name and their corresponding HandleFunc
// the root handler "/" must remain at the end of the function
func (app *App) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	app.ENDPOINTS = Endpoints{
		"GET /about":                        app.HndlAbt,
		"GET /health":                       app.HndlHealth,
		"GET /dbhealth":                     app.HndlDBHealth,
		"GET /bronto":                       app.HndlBronto,
		"GET /bball":                        app.HndlBBall,
		"GET /bball/about":                  app.HndlBBallAbt,
		"GET /bball/seasons":                app.HndlSeasons,
		"GET /bball/teams":                  app.HndlTeams,
		"GET /bball/player":                 app.HndlPlayer,
		"GET /bball/games/recent":           app.HndlRecentGames,
		"GET /bball/league/scoring-leaders": app.HndlTopLgPlayers,
		"GET /bball/teamrecs":               app.HndlTeamRecords,
		"GET /bball/v2/players":             app.HndlPlayerV2,
		"/docs/":                            app.ServeDocs,
		"/js/":                              app.JSNostore,
		"/css/":                             app.CSSNostore,
		"/":                                 app.HndlRoot,
	}

	for pattern, handler := range app.ENDPOINTS {
		mux.HandleFunc(pattern, handler)
	}

	app.Lg.Infof("%d endpoints mounted", len(app.ENDPOINTS))
	return mux
}
