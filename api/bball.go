package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jdetok/golib/errd"
)

// /player handler
func (app *App) HndlPlayer(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)

	var rp Resp
	var tId uint64

	team := r.URL.Query().Get("team")
	tId, _ = strconv.ParseUint(team, 10, 64)

	season := r.URL.Query().Get("season")
	player := RemoveDiacritics(r.URL.Query().Get("player"))
	pId, sId := GetpIdsId(app.Players, player, season)

	js, err := rp.GetPlayerDash(app.Database, pId, sId, tId)
	if err != nil {
		msg := fmt.Sprintf("server failed to return player dash for %s", player)
		e.HTTPErr(w, msg, err)
	}
	app.JSONWriter(w, js)
}

// /games/recent handler
func (app *App) HndlRecentGames(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)
	rgs := RecentGames{}
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
				tm.LogoUrl = tm.MakeLogoUrl()
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
