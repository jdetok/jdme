package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

/*
Player struct meant to store basic global data for each player
SeasonIdMax/Min are the player's first and last REGULAR season in their league
PSeasonIdMax/Min are the player's first and last POST SEASON in their league.
these values will default to 0 for players without any recorded games in a past season
*/
type Player struct {
	PlayerId     uint64
	Name         string
	League       string
	SeasonIdMax  uint64
	SeasonIdMin  uint64
	PSeasonIdMax uint64
	PSeasonIdMin uint64
}

/*
primary database query function for the /players endpoint. queries the api
tables in the database sing the passed player, season, team ID to get the
player's stats. defaults to TeamTopScorerDash query, which gets the dash for
the top scorer of the most recent night's games. this is called when the site
loads. the response is scanned into the structs defined in resp.go, before being
marshalled into json and returned to write as the http response
*/
func (r *Resp) GetPlayerDash(db *sql.DB, pId, sId, tId uint64) ([]byte, error) {
	e := errd.InitErr()
	var q string
	var p uint64

	// if 0 is passed as tId, query by player_id. otherwise, query by team_id
	switch tId {
	case 0:
		logd.Logc(fmt.Sprintf("querying player_id: %d | season_id: %d", pId, sId))
		q = pgdb.PlayerDash
		p = pId
	default:
		logd.Logc(fmt.Sprintf("querying team_id: %d | season_id: %d", tId, sId))
		q = pgdb.TeamTopScorerDash
		p = tId
	}

	// query player, scan to structs, call struct functions
	// appends RespObj to r.Results
	if err := r.BuildPlayerRespStructs(db, q, p, sId); err != nil {
		e.Msg = fmt.Sprintf("failed to query playerId %d seasonId %d", p, sId)
		return nil, e.BuildErr(err)
	}

	// marshall Resp struct to JSON, return as []byte
	js, err := json.Marshal(r)
	if err != nil {
		e.Msg = "failed to marshal structs to json"
		return nil, e.BuildErr(err)
	}
	return js, nil
}

// query player, scan to structs, call struct functions
// appends RespObj to r.Results
// separated from GetPlayerDash 09/24/2025
func (r *Resp) BuildPlayerRespStructs(db *sql.DB, q string, p, sId uint64) error {
	e := errd.InitErr()

	// QUERY SEASON PLAYERDASH FOR pId OR FOR TOP SCORER OF TEAM (tId) PASSED
	rows, err := db.Query(q, p, sId)
	if err != nil {
		e.Msg = "error during player dash query"
		return e.BuildErr(err)
	}

	var t RespSeasonTmp // temp seasons for NBA/WNBA, handled after loop
	var rp RespObj
	for rows.Next() {
		// temp structs, handled in hndlRespRow
		var s RespPlayerStats
		var p RespPlayerSznOvw
		// 8/6 2PM - MOVED Season/WSeason FROM END TO AFTER SeasonId
		rows.Scan( // MUST BE IN ORDER OF QUERY
			&rp.Meta.PlayerId, &rp.Meta.TeamId, &rp.Meta.League,
			&rp.Meta.SeasonId, &t.Season, &t.WSeason, &rp.Meta.StatType,
			&rp.Meta.Player, &rp.Meta.Team, &rp.Meta.TeamName,
			&rp.SeasonOvw.GamesPlayed, &p.Minutes,
			&s.Box.Points, &s.Box.Assists, &s.Box.Rebounds,
			&s.Box.Steals, &s.Box.Blocks,
			&s.Shtg.Fg.Makes, &s.Shtg.Fg.Attempts, &s.Shtg.Fg.Percent,
			&s.Shtg.Fg3.Makes, &s.Shtg.Fg3.Attempts, &s.Shtg.Fg3.Percent,
			&s.Shtg.Ft.Makes, &s.Shtg.Ft.Attempts, &s.Shtg.Ft.Percent)
		// switch on stat type to assign stats to appropriate struct
		rp.HandleStatTypeSznOvw(&p, &s)
	}

	// assign nba or wnba season only based on league
	t.SwitchSznByLeague(&rp.Meta.League, &rp.Meta.Season)

	// build table captions & image urls
	rp.Meta.MakePlayerDashCaptions()
	rp.Meta.MakeHeadshotUrl()
	rp.Meta.MakeTeamLogoUrl()

	// append built respObj to Resp, return
	r.Results = append(r.Results, rp)
	return nil
}

/*
switch between totals (sums) and pergame (averages) stats based on the
Meta.StatType field
*/
func (rp *RespObj) HandleStatTypeSznOvw(p *RespPlayerSznOvw, s *RespPlayerStats) {
	switch rp.Meta.StatType {
	case "avg":
		rp.SeasonOvw.MinutsPerGame = p.Minutes
		rp.PerGame.Box = s.Box
		rp.PerGame.Shtg = s.Shtg
	case "tot":
		rp.SeasonOvw.Minutes = p.Minutes
		rp.Totals.Box = s.Box
		rp.Totals.Shtg = s.Shtg
	}
}

/*
accept slice of Player structs and a season id, call slicePlayerSzn to create
a new slice with only players from the specified season. then, generate a
random number and return the player at that index in the slice
*/
func RandomPlayerId(pl []Player, cs *CurrentSeasons, sId uint64, lg string) uint64 {
	players, _ := SlicePlayersSzn(pl, cs, sId, lg)
	numPlayers := len(players)
	randNum := rand.IntN(numPlayers)
	return players[randNum].PlayerId
}

/*
player name and season ID from get request passed here, returns the player's
ID and the season ID. if 'player' variable == "random", the randPlayer function
is called. a player ID also can be passed as the player parameter, it will just
be converted to an int and returned
*/
func ValidatePlayerSzn(players []Player, cs *CurrentSeasons, player string, seasonId string, lg string, errStr *string) (uint64, uint64) {
	sId, _ := strconv.ParseUint(seasonId, 10, 32)
	var pId uint64

	if player == "random" { // call randplayer function
		pId = RandomPlayerId(players, cs, sId, lg)
	} else if _, err := strconv.ParseUint(player, 10, 64); err == nil {
		// if it's numeric keep it and convert to uint64
		pId, _ = strconv.ParseUint(player, 10, 64)
	} else { // search name through players list
		for _, p := range players {
			if p.Name == player { // return match playerid (uint32) as string
				pId = p.PlayerId
			}
		}
	}

	// loop through players to check that queried season is within min-max seasons
	for _, p := range players {
		if p.PlayerId == pId {
			return pId, HandleSeasonId(sId, &p, errStr)
		}
	}
	return pId, sId
}

func (m *RespPlayerMeta) MakePlayerDashCaptions() {
	var delim string = "|"
	m.Caption = fmt.Sprintf("%s %s %s", m.Player, delim, m.TeamName)
	m.CaptionShort = fmt.Sprintf("%s %s %s", m.Player, delim, m.Team)
	m.BoxCapTot = fmt.Sprintf("Box Totals %s %s", delim, m.Season)
	m.BoxCapAvg = fmt.Sprintf("Box Averages %s %s", delim, m.Season)
	m.ShtgCapTot = fmt.Sprintf("Shooting Totals %s %s", delim, m.Season)
	m.ShtgCapAvg = fmt.Sprintf("Shooting Averages %s %s", delim, m.Season)
}
