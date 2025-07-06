package main

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/internal/env"
	"github.com/jdetok/go-api-jdeko.me/internal/errs"
	"github.com/jdetok/go-api-jdeko.me/internal/jsonops"
	"github.com/jdetok/go-api-jdeko.me/internal/logs"
	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
	"github.com/jdetok/go-api-jdeko.me/internal/store"
)

/* TODO
- ENDPOINT: MOST RECENT DATE WITH GAME(S)
-- highest scoring player's id for headshot

- ENDPOINT: CURRENT LEAGUE RECORDS
*/
// SERVE MAIN NBA PAGES
var nbaPath string = (env.GetString("STATIC_PATH") + "/nba.html")
var nbaDevPath string = (env.GetString("STATIC_PATH") + "/dev/nba.html")

func (app *application) bballHandler(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	http.ServeFile(w, r, nbaPath)
}

func (app *application) bballDevHandler(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	http.ServeFile(w, r, nbaDevPath)
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

// RANDOM PLAYER BUTTON
func (app *application) getRandPlayer(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	numPlayers := len(app.players)
	randNum := rand.IntN(numPlayers)

	w.Header().Set("Content-Type", "application/json") 
	json.NewEncoder(w).Encode(map[string]string{
		"playerId": strconv.FormatUint(app.players[randNum].PlayerId, 10),
		"player": app.players[randNum].Name,
		"league": app.players[randNum].League,
	})
	// random number in range of len(players) to return random player
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

// USED TO READ STATS FROM JSON FILES IN getStats FUNC
func respFromFile(w *http.ResponseWriter, f string) []byte {
	e := errs.ErrInfo{Prefix: "json file read"}
	js, err := jsonops.ReadJSON(env.GetString("CACHE_PATH") + f)
	if err != nil {
		e.Msg = ("error reading json file: " + f)
		errs.HTTPErr(*w, e.Error(err))
	}
	return js
}

/* 
HANDLE /bball/players REQUESTS
EX. ALL PLAYERS: ?lg=nba&stype=tot&player=all
EX. SPECIFIC PLAYER: ?lg=nba&stype=avg&player=tyrese%20haliburton
*/
func (app *application) getStats(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)

	lg := r.URL.Query().Get("lg")
	sType := r.URL.Query().Get("stype")
	player := r.URL.Query().Get("player")
	
	switch lg { // OUTER SWTICH: NBA/WNBA	
	case "nba": // NBA SWITCH TOTALS/AVERAGES
		switch sType {
		case "tot": // NBA TOTALS
			switch player {
			case "all": // NBA RS TOTALS ALL PLAYERS
				js := respFromFile(&w, "/nba_rs_totals.json")
				app.JSONWriter(w, js)
			default: // NBA RS TOTALS SPECIFIC PLAYER
				playerId := store.SearchPlayers(app.players, player)
				js := mariadb.SelectLgPlayer(app.database, &w, 
					mariadb.LgPlayerStat.Q, lg, string(playerId))
				app.JSONWriter(w, js)
			}
		case "avg": // NBA AVERAGES
			switch player {
			case "all": // NBA RS AVG ALL PLAYERS
				js := respFromFile(&w, "/nba_rs_avgs.json")
				app.JSONWriter(w, js)
			default:  // NBA RS AVG SPECIFIC PLAYER
				playerId := store.SearchPlayers(app.players, player)
				js := mariadb.SelectLgPlayer(app.database, &w, 
					mariadb.LgPlayerAvg.Q, lg, string(playerId))
				app.JSONWriter(w, js)
			}
		}
	case "wnba":
		switch sType { // WNBA SWITCH TOTALS/AVERAGES
		case "tot": // WNBA TOTALS
			switch player {  // ALL WNBA TOTALS
			case "all":
				js := respFromFile(&w, "/wnba_rs_totals.json")
				app.JSONWriter(w, js)
			default: // SPECIFIC WNBA PLAYER TOTALS
				playerId := store.SearchPlayers(app.players, player)
				js := mariadb.SelectLgPlayer(app.database, &w, 
					mariadb.LgPlayerStat.Q, lg, string(playerId))
				app.JSONWriter(w, js)
			}
		case "avg": // WNBA AVERAGES
			switch player {
			case "all": // ALL WNBA AVERAGES
				js := respFromFile(&w, "/wnba_rs_avgs.json")
				app.JSONWriter(w, js)
			default: // SPECIFIC WNBA PLAYER AVERAGES
				playerId := store.SearchPlayers(app.players, player)
				js := mariadb.SelectLgPlayer(app.database, &w, 
					mariadb.LgPlayerAvg.Q, lg, string(playerId))
				app.JSONWriter(w, js)
			}
		}
	}
}