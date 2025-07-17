package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/api/cache"
	"github.com/jdetok/go-api-jdeko.me/api/resp"
	"github.com/jdetok/go-api-jdeko.me/applog"
)

func (app *application) getPlayerDash(w http.ResponseWriter, r *http.Request) {
	e := applog.AppErr{Process: "player dash endpoint", IsHTTP: true}
	applog.LogHTTP(r)
	var rp resp.Resp

	var tId uint64
	team := r.URL.Query().Get("team")
	tId, _ = strconv.ParseUint(team, 10, 64)

	season := r.URL.Query().Get("season")
	player := cache.Unaccent(r.URL.Query().Get("player"))
	pId, sId := resp.GetpIdsId(app.players, player, season)

	js, err := rp.GetPlayerDash(app.database, pId, sId, tId)
	if err != nil {
		e.Msg = "failed to get player dash"
		e.MsgHTTP = fmt.Sprintf("server failed to return player dash for %s", player)
		e.HTTPErr(w, e.BuildError(err))
	}
	app.JSONWriter(w, js)
}

func (app *application) getGamesRecentNew(w http.ResponseWriter, r *http.Request) {
	e := applog.AppErr{Process: "recent games endpoint"}
	applog.LogHTTP(r)
	rgs := resp.RecentGames{}
	js, err := rgs.GetRecentGames(app.database)
	if err != nil {
		e.Msg = "failed to get recent games"
		e.MsgHTTP = "server failed to return recent games"
		e.HTTPErr(w, e.BuildError(err))
	}
	app.JSONWriter(w, js)
}

// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *application) getSeasons(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	season := r.URL.Query().Get("szn")
	w.Header().Set("Content-Type", "application/json")
	if season == "" {
		json.NewEncoder(w).Encode(app.seasons)
	} else {
		for _, szn := range app.seasons {
			if season == szn.SeasonId {
				json.NewEncoder(w).Encode(map[string]string{
					"szn": season,
				})
			}
		}
	}
}

// FOR TEAMS SELECTOR - CALLED ON PAGE LOAD
func (app *application) getTeams(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	team := r.URL.Query().Get("team")
	w.Header().Set("Content-Type", "application/json")
	if team == "" {
		json.NewEncoder(w).Encode(app.teams)
	} else {
		for _, tm := range app.teams {
			if team == tm.TeamAbbr {
				tm.LogoUrl = tm.MakeLogoUrl()
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
