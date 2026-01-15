package api

import (
	"net/http"
	"strconv"

	"github.com/jdetok/jdme/pkg/memd"
)

// recent games, team records, top players, etc

// /games/recent handler
func (app *App) HndlRecentGames(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)

	var rgs memd.RecentGames
	js, err := rgs.GetRecentGames(app.DB)
	if err != nil {
		msg := "server failed to return recent games"
		app.Lg.HTTPErr(w, err, http.StatusUnprocessableEntity, msg)
	}
	app.JSONWriter(w, js)
}

// top numPl players for each league
func (app *App) HndlTopLgPlayers(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)

	// new LgTopPlayers to get from in memory store
	var lt memd.LgTopPlayers

	// get number of players from query string & convert to int
	numPlStr := r.URL.Query().Get("num")
	numPl, err := strconv.ParseUint(numPlStr, 10, 64)
	if err != nil {
		msg := "failed to convert numPlStr to int"
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}

	// append numPl to NBA/WNBA LgTopPlayer from memory store
	for i := range numPl {
		lt.NBATop = append(lt.NBATop, app.Store.TopLgPlayers.NBATop[i])
		lt.WNBATop = append(lt.WNBATop, app.Store.TopLgPlayers.WNBATop[i])
	}

	// marshal new LgTopPlayers to json
	js, err := memd.MarshalTopPlayers(&lt)
	if err != nil {
		msg := "failed to marshal top 5 league players struct to JSON"
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	app.JSONWriter(w, js)
}

// team records for current/most recent reg. seasons for each league
func (app *App) HndlTeamRecords(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)

	js, err := memd.TeamRecordsJSON(&app.Store.TeamRecs)
	if err != nil {
		msg := "failed to marshal team records struct to JSON"
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	app.JSONWriter(w, js)
}
