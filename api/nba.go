package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/internal/env"
	"github.com/jdetok/go-api-jdeko.me/internal/errs"
	"github.com/jdetok/go-api-jdeko.me/internal/logs"
	"github.com/jdetok/go-api-jdeko.me/internal/store"
)

var nbaDevPath string = (env.GetString("DEV_PATH") + "/bball/nba.html")
var bballPath string = (env.GetString("BBALL_PATH") + "/nba.html")

func (app *application) bballHandler(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	http.ServeFile(w, r, bballPath)
}

func (app *application) bballDevHandler(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	http.ServeFile(w, r, nbaDevPath)
}

func (app *application) getPlayerDash(w http.ResponseWriter, r *http.Request) {
	// player & season are params
	e := errs.ErrInfo{Prefix: "player dash endpoint"}
	logs.LogHTTP(r)

	season := r.URL.Query().Get("season")
	player := store.Unaccent(r.URL.Query().Get("player"))
	team := r.URL.Query().Get("team")
	var tId uint64
	tId, _ = strconv.ParseUint(team, 10, 64)

	pId, sId := store.GetpIdsId(app.players, player, season)

	var rp store.Resp
	js, err := rp.GetPlayerDash(app.database, pId, sId, tId)
	// js, err := rp.GetPlayerDash(app.database, pId, sId, tId)
	if err != nil {
		e.Msg = "failed to get player dash"
		errs.HTTPErr(w, e.Error(err))
	}
	app.JSONWriter(w, js)
}

func (app *application) getGamesRecentNew(w http.ResponseWriter, r *http.Request) {
	e := errs.ErrInfo{Prefix: "recent games endpoint"}
	logs.LogHTTP(r)
	// js, err := mariadb.DBJSONResposne(app.database, mariadb.RecentGamePlayers.Q)
	rgs := store.RecentGames{}
	js, err := rgs.GetRecentGames(app.database)
	if err != nil {
		e.Msg = "failed to get games"
		errs.HTTPErr(w, e.Error(err))
	}
	app.JSONWriter(w, js)
}

func (app *application) getTopScorerNew(w http.ResponseWriter, r *http.Request) {
	e := errs.ErrInfo{Prefix: "top scorers endpoint"}
	logs.LogHTTP(r)
	ts := store.TopScorers{}
	js, err := ts.GetTopScorers(app.database)
	if err != nil {
		e.Msg = ("failed to get top scorers")
		errs.HTTPErr(w, e.Error(err))
	}
	app.JSONWriter(w, js)
}

// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *application) getSeasons(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
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
	logs.LogHTTP(r)
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

// CALLED BY JS TO GET PLAYER'S HEADSHOT (DEPRECATE, USE ONE CALL)
func (app *application) getPlayerId(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	player := r.URL.Query().Get("player")
	logs.LogDebug("Player Requested: " + player)

	playerId := store.SearchPlayers(app.players, player)
	logs.LogDebug("PlayerId Return: " + playerId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"playerId": playerId,
	})
}
