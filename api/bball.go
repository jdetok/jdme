package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/api/resp"
	"github.com/jdetok/go-api-jdeko.me/api/store"
	"github.com/jdetok/golib/errd"
)

func (app *App) playerDashHndl(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)

	var rp resp.Resp
	var tId uint64

	team := r.URL.Query().Get("team")
	tId, _ = strconv.ParseUint(team, 10, 64)

	season := r.URL.Query().Get("season")
	player := store.Unaccent(r.URL.Query().Get("player"))
	pId, sId := resp.GetpIdsId(app.players, player, season)

	js, err := rp.GetPlayerDash(app.database, pId, sId, tId)
	if err != nil {
		msg := fmt.Sprintf("server failed to return player dash for %s", player)
		e.HTTPErr(w, msg, err)
	}
	app.JSONWriter(w, js)
}

// come back to this - used in top scorer maybe?
func (app *App) recGameHndl(w http.ResponseWriter, r *http.Request) {
	e := errd.InitErr()
	LogHTTP(r)
	rgs := resp.RecentGames{}
	js, err := rgs.GetRecentGames(app.database)
	if err != nil {
		e.Msg = "failed to get recent games"
		msg := "server failed to return recent games"
		e.HTTPErr(w, msg, err)
	}
	app.JSONWriter(w, js)
}

// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *App) seasonsHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	season := r.URL.Query().Get("szn")
	w.Header().Set("Content-Type", "App/json")
	if season == "" { // send all szns when szn is not in q str, used most often
		json.NewEncoder(w).Encode(app.seasons)
	} else {
		for _, szn := range app.seasons {
			if season == szn.SeasonId { // validate szn from q string
				json.NewEncoder(w).Encode(map[string]string{
					"szn": season,
				})
			}
		}
	}
}

// FOR TEAMS SELECTOR - CALLED ON PAGE LOAD
func (app *App) teamsHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	team := r.URL.Query().Get("team")
	w.Header().Set("Content-Type", "App/json")
	if team == "" { // send all teams when team is not in q str, used most often
		json.NewEncoder(w).Encode(app.teams)
	} else { // read & valid team from q string, not yet used 8/6
		for _, tm := range app.teams {
			if team == tm.TeamAbbr {
				tm.LogoUrl = tm.MakeLogoUrl()
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
