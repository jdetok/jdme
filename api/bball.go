package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PlayerQuery struct {
	Player string
	Team   string
	Season string
	League string
}

type PlayerTeamSeason struct {
	Player string
	Team   string
	Season string
}

type PQueryIds struct {
	PId uint64
	TId uint64
	SId uint64
}

// top numPl players for each league
func (app *App) HndlTopLgPlayers(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)

	// new LgTopPlayers to get from in memory store
	var lt LgTopPlayers

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
	js, err := MarshalTopPlayers(&lt)
	if err != nil {
		msg := "failed to marshal top 5 league players struct to JSON"
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	app.JSONWriter(w, js)
}

// team records for current/most recent reg. seasons for each league
func (app *App) HndlTeamRecords(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)

	js, err := TeamRecordsJSON(&app.Store.TeamRecs)
	if err != nil {
		msg := "failed to marshal team records struct to JSON"
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	app.JSONWriter(w, js)
}

// /player handler
/*
- create PlayerQuery pq struct to hold query string params as strings
- create PQueryIds iq to hold them as uint64s
- get all as strings, remove accents from player names
- call ValidatePlayerSzn
	- pass pq strings, get iq uints back
	- attempts to convert teamId and seasonId to ints
	- calls RandomPlayerId if player name is "random"
	- if passed as a player id it assigns that to iq.PId and moves on
	- otherwise searches passed player name again the in-mem player slice
	- finally calls HandleSeasonId to validate the seasonId
- call GetPlayerDash
	- pass iq (result of ValidatePlayerSzn)
*/
func (app *App) HndlPlayer(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)

	// pq holds all query parameters as strings
	var pq PlayerQuery
	var iq PQueryIds
	// var iq PQueryIds
	var rp Resp
	// var tId uint64

	// get and convert team from query string
	pq.Team = r.URL.Query().Get("team")

	// get season from query string
	pq.Season = r.URL.Query().Get("season")

	//get league from query string
	pq.League = r.URL.Query().Get("league")

	// get player from query string
	pq.Player = RemoveDiacritics(r.URL.Query().Get("player"))

	// validate player & get playerid/season id
	iq, err := ValidatePlayerSzn(app.Store.Players, &app.Store.CurrentSzns, &pq, &rp.ErrorMsg)
	if err != nil {
		msg := fmt.Sprintf("validate player %s", pq.Player)
		app.Lg.HTTPErr(w, err, http.StatusUnprocessableEntity, msg)
	}

	// query the player & build JSON response, returned as []byte to write
	js, err := rp.GetPlayerDash(app.Database, &iq)
	if err != nil {
		msg := fmt.Sprintf("server failed to return player dash for %s", pq.Player)
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	// write JSON response
	app.JSONWriter(w, js)
}

// /games/recent handler
func (app *App) HndlRecentGames(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)

	var rgs RecentGames
	js, err := rgs.GetRecentGames(app.Database)
	if err != nil {
		msg := "server failed to return recent games"
		app.Lg.HTTPErr(w, err, http.StatusUnprocessableEntity, msg)
	}
	app.JSONWriter(w, js)
}

// /seasons handler
// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *App) HndlSeasons(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)
	season := r.URL.Query().Get("szn")
	w.Header().Set("Content-Type", "application/json")
	if season == "" { // send all szns when szn is not in q str, used most often
		json.NewEncoder(w).Encode(app.Store.Seasons)
	} else {
		for _, szn := range app.Store.Seasons {
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
	app.Lg.LogHTTP(r)
	team := r.URL.Query().Get("team")
	w.Header().Set("Content-Type", "application/json")
	if team == "" { // send all teams when team is not in q str, used most often
		json.NewEncoder(w).Encode(app.Store.Teams)
	} else { // read & valid team from q string, not yet used 8/6
		for _, tm := range app.Store.Teams {
			if team == tm.TeamAbbr {

				tm.LogoUrl = MakeTeamLogoUrl(tm.League, tm.TeamId)
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
