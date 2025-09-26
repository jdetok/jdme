package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// TODO: if sId 88888 doesn't need to do season
func VerifyPlayerTeam(db *sql.DB, iq *PQueryIds) (bool, error) {
	e := errd.InitErr()
	fmt.Println("VerifyPlayerTeam func called")
	rows, err := db.Query(pgdb.PlayerTeamBool, iq.PId, iq.TId)
	if err != nil {
		e.Msg = fmt.Sprintf("error verifying season|team|player | %d | %d | %d",
			iq.SId, iq.TId, iq.PId)
		return false, e.BuildErr(err)
	}
	// just need to know if a row exists -
	return rows.Next(), nil
}

// TODO: if sId 88888 doesn't need to do season
func VerifyPlayerTeamSeason(db *sql.DB, iq *PQueryIds) (bool, error) {
	e := errd.InitErr()
	fmt.Println("VerifyPlayerTeamSeason func called")
	rows, err := db.Query(pgdb.VerifyTeamSzn, iq.SId, iq.TId, iq.PId)
	if err != nil {
		e.Msg = fmt.Sprintf("error verifying season|team|player | %d | %d | %d",
			iq.SId, iq.TId, iq.PId)
		return false, e.BuildErr(err)
	}
	// just need to know if a row exists -
	return rows.Next(), nil
}

/*
primary database query function for the /players endpoint. queries the api
tables in the database sing the passed player, season, team ID to get the
player's stats. defaults to TeamTopScorerDash query, which gets the dash for
the top scorer of the most recent night's games. this is called when the site
loads. the response is scanned into the structs defined in resp.go, before being
marshalled into json and returned to write as the http response
*/
func (r *Resp) GetPlayerDash(db *sql.DB, iq *PQueryIds) ([]byte, error) {
	e := errd.InitErr()
	// query player, scan to structs, call struct functions
	// appends RespObj to r.Results
	if err := r.BuildPlayerRespStructs(db, iq); err != nil {
		e.Msg = fmt.Sprintf("failed to query playerId %d seasonId %d", iq.PId, iq.SId)
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
/*
swtiches query and arguments based on whether teamId = 0
*/
func (r *Resp) BuildPlayerRespStructs(db *sql.DB, iq *PQueryIds) error {
	e := errd.InitErr()
	var args []uint64
	var q string
	// QUERY SEASON PLAYERDASH FOR pId OR FOR TOP SCORER OF TEAM (tId) PASSED
	if iq.TId > 0 {
		var ptValid bool
		var err error
		if iq.SId == 88888 {
			ptValid, err = VerifyPlayerTeam(db, iq)
			if err != nil {
				e.Msg = "error executing team player validation query"
				return e.BuildErr(err)
			}
		} else {
			ptValid, err = VerifyPlayerTeamSeason(db, iq)
			if err != nil {
				e.Msg = "error executing team player validation query"
				return e.BuildErr(err)
			}
		}

		if !(ptValid) {
			errmsg := fmt.Sprintf(
				"%d did not play for %d in %d",
				iq.PId, iq.TId, iq.SId)
			r.ErrorMsg = errmsg
			e.Msg = errmsg
			// return e.NewErr()
		}
		fmt.Printf(
			"Player %d Team %d Season %d validated in BuildPlayerRespStructs func\n",
			iq.PId, iq.TId, iq.SId)
		args = []uint64{iq.PId, iq.TId}
		q = pgdb.TstTeamPlayer
	} else {
		args = []uint64{iq.PId, iq.SId}
		q = pgdb.PlayerDash
	}
	// rows , err := db.Query(q, iq.PId, iq.SId)
	rows, err := db.Query(q, args[0], args[1])
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
		// two rows are returned, one for each stattype
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
	// rp.Meta.
	rp.Meta.TeamLogoUrl = MakeTeamLogoUrl(rp.Meta.League, strconv.FormatUint(rp.Meta.TeamId, 10))

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
func RandomPlayerId(pl []Player, cs *CurrentSeasons, pq *PlayerQuery) uint64 {
	players, _ := SlicePlayersSzn(pl, cs, pq)
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
func ValidatePlayerSzn(
	players []Player,
	cs *CurrentSeasons,
	pq *PlayerQuery,
	errStr *string) (PQueryIds, error) {
	//
	e := errd.InitErr()
	var iq PQueryIds

	sId, err := strconv.ParseUint(pq.Season, 10, 32)
	if err != nil {
		e.Msg = fmt.Sprintf("failed to convert season %s to int", pq.Season)
		return iq, e.BuildErr(err)
	}

	tId, err := strconv.ParseUint(pq.Team, 10, 64)
	if err != nil {
		e.Msg = fmt.Sprintf("failed to convert team %s to int", pq.Team)
		return iq, e.BuildErr(err)
	}
	iq.TId = tId
	// var pId uint64

	if pq.Player == "random" { // call randplayer function
		iq.PId = RandomPlayerId(players, cs, pq)
	} else {
		// check if it can be converted to number
		pId, err := strconv.ParseUint(pq.Player, 10, 64)
		if err == nil {
			// if it's numeric keep it and convert to uint64
			iq.PId = pId
		} else { // search name through players list
			for _, p := range players {
				if p.Name == pq.Player { // return match playerid (uint32) as string
					iq.PId = p.PlayerId
				}
			}
		}
	}

	// loop through players to check that queried season is within min-max seasons
	for _, p := range players {
		if p.PlayerId == iq.PId {
			if iq.TId > 0 {
				iq.SId = HandleSeasonId(sId, &p, true, errStr)
			} else {
				iq.SId = HandleSeasonId(sId, &p, false, errStr)
			}
			return iq, nil
		}
	}
	e.Msg = "player not in memory"
	return iq, e.NewErr()
}

// fill RespPlayerMeta captions fields with formatted strings for each playerdash
// table's caption
func (m *RespPlayerMeta) MakePlayerDashCaptions() {
	var delim string = "|"
	m.Caption = fmt.Sprintf("%s %s %s", m.Player, delim, m.TeamName)
	m.CaptionShort = fmt.Sprintf("%s %s %s", m.Player, delim, m.Team)
	m.BoxCapTot = fmt.Sprintf("Box Totals %s %s", delim, m.Season)
	m.BoxCapAvg = fmt.Sprintf("Box Averages %s %s", delim, m.Season)
	m.ShtgCapTot = fmt.Sprintf("Shooting Totals %s %s", delim, m.Season)
	m.ShtgCapAvg = fmt.Sprintf("Shooting Averages %s %s", delim, m.Season)
}

/*
use the transform package to remove accidentals
e.g. Dončić becomes doncic
*/
func RemoveDiacritics(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	output, _, _ := transform.String(t, input)
	return output
}

// use league and player id to build the URL containing a player's headshot
func (m *RespPlayerMeta) MakeHeadshotUrl() {
	lg := strings.ToLower(m.League)
	pId := strconv.Itoa(int(m.PlayerId))
	m.HeadshotUrl = fmt.Sprintf(
		`https://cdn.%s.com/headshots/%s/latest/1040x760/%s.png`,
		lg, lg, pId)
}

// use league and team id to build team logo URLs
func MakeTeamLogoUrl(league, teamId string) string {
	lg := strings.ToLower(league)
	// tId := strconv.Itoa(int(teamId))
	return fmt.Sprintf(
		`https://cdn.%s.com/logos/%s/%s/primary/L/logo.svg`,
		lg, lg, teamId)
}

/*
accept a season id and a pointer to a Player struct, validate the player was active
in the passed season, return a valid season ID if not. if season id starts with an
8 the player's max regular season will be returned. if it starts with a 7, their
max playoff season will be returned. if it starts with a 4, it will first verify
player has played in a playoff game, and will return their max regular season if
they haven't. a season id starting with 2 will return a regular season. for both
regular season and playoffs, the function will verify the player played in said
season, and return either their max or min (whichever is closer) season  if they
did not
*/
func HandleSeasonId(sId uint64, p *Player, team bool, errStr *string) uint64 {
	if sId == 99999 || sId == 29999 { // agg seasons
		msg := fmt.Sprintf("aggregate season requested%d | %d\n", sId, sId)
		fmt.Println(msg)
		return sId
	} else if sId == 88888 {
		if team {
			return sId
		}
		msg := fmt.Sprintf("returning latest regular season for player%d | %d\n",
			sId, p.SeasonIdMax)
		fmt.Println(msg)
		return p.SeasonIdMax // return most recent season
	} else if sId >= 40000 && sId < 50000 {
		if p.PSeasonIdMax < 40000 { // player has no playeroff, return max reg season
			msg := fmt.Sprintf(
				"%s has not played in the post-season | displaying latest regular season stats",
				p.Name)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMax // return reg season if player has no playoffs
		}
		if sId == 49999 {
			msg := fmt.Sprintf(
				"requested career playoff stats %d | %d\n",
				sId, sId)
			// *errStr = msg
			fmt.Println(msg)
			return sId
		}
		if sId > p.PSeasonIdMax {
			msg := fmt.Sprintf(
				// "szn > playoff max, returning playoff max%d | %d\n",
				// sId, p.PSeasonIdMax)
				"%d was after %s's last playoff season | displaying the %d playoffs",
				sId, p.Name, p.PSeasonIdMax)
			*errStr = msg
			fmt.Println(msg)
			return p.PSeasonIdMax
		}
		if sId < p.PSeasonIdMin {
			msg := fmt.Sprintf(
				"the first playoffs for %s was the %d season",
				p.Name, p.PSeasonIdMin)
			*errStr = msg
			fmt.Println(msg)
			return p.PSeasonIdMin
		}
	} else if sId >= 20000 && sId < 30000 {
		if sId > p.SeasonIdMax {
			msg := fmt.Sprintf(
				"%s has not played games in the %d season | displaying %d stats instead\n",
				p.Name, sId, p.SeasonIdMax)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMax
		}
		if sId < p.SeasonIdMin {
			msg := fmt.Sprintf(
				"%s was not in the league yet for the %d season | displaying their rookie season %d stats instead\n",
				p.Name, sId, p.SeasonIdMin)
			*errStr = msg
			fmt.Println(msg)
			return p.SeasonIdMin
		}
	}
	msg := fmt.Sprintf("validated: %d | %d\n", sId, sId)
	fmt.Println(msg)
	return sId
}

// tId, err := strconv.ParseUint(pq.Team, 10, 64)
// iq.TId = tId
// if err != nil {
// 	msg := fmt.Sprintf("error converting %v to int", pq.Team)
// 	e.HTTPErr(w, msg, err)
// }

// func SlicePlayerTeam(players *[]Player, teamId string) ([]Player, error) {
// 	e := errd.InitErr()
// 	tId, err := strconv.ParseUint(teamId, 10, 64)
// 	if err != nil {
// 		e.Msg = fmt.Sprintf("error converting %s to uint64")
// 		return nil, e.BuildErr(err)
// 	}

// 	for _, p := range *players {
// 		if p.
// 	}

// }

/*
accept the slice of all players and a seasonId, return a slice with just the
active players from the passed season id
*/
// func SlicePlayersSzn(players []Player, seasonId uint64) ([]Player, error) {
func SlicePlayersSzn(players []Player, cs *CurrentSeasons, pq *PlayerQuery) ([]Player, error) {
	e := errd.InitErr()
	var plslice []Player

	if pq.Team != "0" {
		fmt.Println("team not 0")
	}

	seasonId, err := strconv.ParseUint(pq.Season, 10, 64)
	if err != nil {
		e.Msg = fmt.Sprintf("failed to convert %s to int", pq.Season)
		return nil, e.BuildErr(err)
	}
	//
	// get struct with current seasons
	sl := cs.LgSznsByMonth(time.Now())

	for _, p := range players { // EXPAND THIS IF TO CATCH PLAYOFF SEASONS AS WELL

		// handle random season id
		if seasonId == 88888 {
			switch p.League {
			case "nba":
				seasonId = sl.SznId
			case "wnba":
				seasonId = sl.WSznId
			}
		}

		if seasonId == 49999 {
			if p.PSeasonIdMin > 0 && p.SeasonIdMax >= (sl.WSznId-3) {
				if pq.League == "all" || pq.League == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		if seasonId == 29999 {
			if p.SeasonIdMax >= (sl.WSznId - 3) {
				if pq.League == "all" || pq.League == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		// append players to the random slice if the passed season id between player min and max season
		if seasonId >= 20000 && seasonId < 30000 {
			if seasonId <= p.SeasonIdMax && seasonId >= p.SeasonIdMin {
				if pq.League == "all" || pq.League == p.League {
					plslice = append(plslice, p)
				}
			}
		}

		if seasonId >= 40000 && seasonId < 50000 {
			if seasonId <= p.PSeasonIdMax && seasonId >= p.PSeasonIdMin {
				if pq.League == "all" || pq.League == p.League {
					plslice = append(plslice, p)
				}
			}
		}

	}
	return plslice, nil
}

// accept pointers of league and season, switch season/wseason on league
func (t *RespSeasonTmp) SwitchSznByLeague(league *string, season *string) {
	switch *league {
	case "NBA":
		*season = t.Season
	case "WNBA":
		*season = t.WSeason
	}
}
