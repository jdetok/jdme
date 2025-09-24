package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jdetok/golib/errd"
)

func (app *App) HndlTopLgPlayers(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)

	numPl := r.URL.Query().Get("num")

	lt, err := QueryTopLgPlayers(app.Database, &app.CurrentSzns, numPl)
	if err != nil {
		msg := "failed to query top 5 league players"
		e.HTTPErr(w, msg, err)
	}
	js, err := MarshalTop5(&lt)
	if err != nil {
		msg := "failed to marshal top 5 league players struct to JSON"
		e.HTTPErr(w, msg, err)
	}
	app.JSONWriter(w, js)
}

// /player handler
func (app *App) HndlPlayer(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)

	var rp Resp
	var tId uint64

	// get and convert team from query string
	team := r.URL.Query().Get("team")
	tId, err := strconv.ParseUint(team, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("error converting %v to int", team)
		e.HTTPErr(w, msg, err)
	}

	// get season from query string
	season := r.URL.Query().Get("season")

	//get league from query string
	lg := r.URL.Query().Get("league")

	// get player from query string
	player := RemoveDiacritics(r.URL.Query().Get("player"))

	// validate player & get playerid/season id
	pId, sId := ValidatePlayerSzn(app.Players, &app.CurrentSzns, player, season, lg, &rp.ErrorMsg)

	// query the player & build JSON response, returned as []byte to write
	js, err := rp.GetPlayerDash(app.Database, pId, sId, tId)
	if err != nil {
		msg := fmt.Sprintf("server failed to return player dash for %s", player)
		e.HTTPErr(w, msg, err)
	}

	// write JSON response
	app.JSONWriter(w, js)
}

// /games/recent handler
func (app *App) HndlRecentGames(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)
	// rgs := RecentGames{}
	var rgs RecentGames
	js, err := rgs.GetRecentGames(app.Database)
	if err != nil {
		e.Msg = "failed to get recent games"
		msg := "server failed to return recent games"
		e.HTTPErr(w, msg, err)
	}
	app.JSONWriter(w, js)
}

// /seasons handler
// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *App) HndlSeasons(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	season := r.URL.Query().Get("szn")
	w.Header().Set("Content-Type", "application/json")
	if season == "" { // send all szns when szn is not in q str, used most often
		json.NewEncoder(w).Encode(app.Seasons)
	} else {
		for _, szn := range app.Seasons {
			if season == szn.SeasonId { // validate szn from q string
				json.NewEncoder(w).Encode(map[string]string{
					"szn": season,
				})
			}
		}
	}
}

// /teams handler
// FOR TEAMS SELECTOR - CALLED ON PAGE LOAD
func (app *App) HndlTeams(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	team := r.URL.Query().Get("team")
	w.Header().Set("Content-Type", "application/json")
	if team == "" { // send all teams when team is not in q str, used most often
		json.NewEncoder(w).Encode(app.Teams)
	} else { // read & valid team from q string, not yet used 8/6
		for _, tm := range app.Teams {
			if team == tm.TeamAbbr {
				tm.LogoUrl = tm.MakeTeamLogoUrl()
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
