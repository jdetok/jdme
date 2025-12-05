package memd

import (
	"database/sql"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

type InMemStore struct {
	Players      []Player
	Seasons      []Season
	Teams        []Team
	CurrentSzns  CurrentSeasons
	TeamRecs     TeamRecords
	TopLgPlayers LgTopPlayers
}

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
	SeasonIdMax  int
	SeasonIdMin  int
	PSeasonIdMax int
	PSeasonIdMin int
}

type Season struct {
	SeasonId string `json:"season_id"`
	Season   string `json:"season"`
	WSeason  string `json:"wseason"`
}

type Team struct {
	League   string `json:"league"`
	TeamId   string `json:"team_id"`
	TeamAbbr string `json:"team"`
	CityTeam string `json:"team_long"`
	LogoUrl  string `json:"-"`
}

// // query db for team season records to populate records table
// // moved to store rather than querying for every request
func UpdateTeamRecords(db *sql.DB, cs *CurrentSeasons) (TeamRecords, error) {
	var team_recs TeamRecords

	sl, err := cs.LgSznsByMonth(time.Now())
	if err != nil {
		return team_recs, err
	}
	rows, err := db.Query(pgdb.TeamSznRecords, sl.SznId, sl.WSznId)
	if err != nil {
		return team_recs, err
	}
	for rows.Next() {
		var tr TeamRecord
		rows.Scan(&tr.League, &tr.SeasonId, &tr.Season, &tr.SeasonDesc,
			&tr.TeamId, &tr.Team, &tr.TeamLong, &tr.Wins, &tr.Losses)

		// append appropriate records based on season
		if tr.League == "NBA" && tr.SeasonId == sl.SznId {
			team_recs.NBARecords = append(team_recs.NBARecords, tr)
		} else if tr.League == "WNBA" && tr.SeasonId == sl.WSznId {
			team_recs.WNBARecords = append(team_recs.WNBARecords, tr)
		}
	}
	return team_recs, nil
}

/*
query the database to update global slice of player structs (in memory player store)
query also gets player's min and max seasons (reg season and playoffs)
*/
func UpdatePlayers(db *sql.DB) ([]Player, error) {
	rows, err := db.Query(pgdb.PlayersSeason)
	if err != nil {
		return nil, err
	}
	var players []Player
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax,
			&p.SeasonIdMin, &p.PSeasonIdMax, &p.PSeasonIdMin)

		// remove diacritics from names in database for imrpvoved searching
		p.Name = clnd.ConvToASCII(p.Name)
		players = append(players, p)
	}
	return players, nil
}

/*
query the database for all seasons, populates global seasons store
example: seasonId: 22025 | season: 2024-25 | WSeason: 2025-26
*/
func UpdateSeasons(db *sql.DB) ([]Season, error) {
	rows, err := db.Query(pgdb.AllSeasons)
	if err != nil {
		return nil, err
	}

	var seasons []Season
	for rows.Next() {
		var szn Season
		rows.Scan(&szn.SeasonId, &szn.Season, &szn.WSeason)
		seasons = append(seasons, szn)
	}
	return seasons, nil
}

// query database for global teams store
func UpdateTeams(db *sql.DB) ([]Team, error) {
	rows, err := db.Query(pgdb.Teams)
	if err != nil {
		return nil, err
	}

	var teams []Team
	for rows.Next() {
		var tm Team
		rows.Scan(&tm.League, &tm.TeamId, &tm.TeamAbbr, &tm.CityTeam)
		teams = append(teams, tm)
	}
	return teams, nil
}
