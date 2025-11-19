package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
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

/*
launch goroutines to update the global players, seasons, and team stores
updates at every interval
*/
func CheckInMemStructs(app *App, interval, threshold time.Duration) {
	e := errd.InitErr()

	// call update func on intial run
	if app.Started == 0 {
		UpdateStructs(app, &e)
		app.Started = 1
	}

	// update structs every interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(app.LastUpdate) > threshold {
			fmt.Printf("refreshing store at %v\n", time.Now().Format(
				"2006-01-02 15:04:05"),
			)
			// calculate current NBA/WNBA seasons
			app.Store.CurrentSzns.GetCurrentSzns(time.Now(), &e)

			// update in memory structs
			UpdateStructs(app, &e)
		}
	}
}

// update players, seasons, and teams in memory structs slices
func UpdateStructs(app *App, e *errd.Err) {
	var err error

	// update in memory players slice
	app.Store.Players, err = UpdatePlayers(app.Database)
	if err != nil {
		e.Msg = "failed to get players"
		fmt.Println(e.BuildErr(err))
	}

	// update in memory seasons slice
	app.Store.Seasons, err = UpdateSeasons(app.Database)
	if err != nil {
		e.Msg = "failed to get seasons"
		fmt.Println(e.BuildErr(err))
	}

	// update in memory teams slice
	app.Store.Teams, err = UpdateTeams(app.Database)
	if err != nil {
		e.Msg = "failed to get teams"
		fmt.Println(e.BuildErr(err))
	}

	// update team records
	app.Store.TeamRecs, err = UpdateTeamRecords(app.Database, &app.Store.CurrentSzns)
	if err != nil {
		e.Msg = "failed to update team records"
		fmt.Println(e.BuildErr(err))
	}

	// update league top players
	app.Store.TopLgPlayers, err = QueryTopLgPlayers(app.Database, &app.Store.CurrentSzns, "20")
	if err != nil {
		e.Msg = "failed to query top league players"
		fmt.Println(e.BuildErr(err))
	}

	// update last update time
	updateTime := time.Now()
	app.LastUpdate = updateTime
	fmt.Printf("finished refreshing store at %v\n", app.LastUpdate)
}

/*
query the database to update global slice of player structs (in memory player store)
query also gets player's min and max seasons (reg season and playoffs)
*/
func UpdatePlayers(db *sql.DB) ([]Player, error) {
	e := errd.InitErr()
	rows, err := db.Query(pgdb.PlayersSeason)
	if err != nil {
		e.Msg = "query failed"
		return nil, e.BuildErr(err)
	}
	var players []Player
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax,
			&p.SeasonIdMin, &p.PSeasonIdMax, &p.PSeasonIdMin)

		// remove diacritics from names in database for imrpvoved searching
		p.Name = RemoveDiacritics(p.Name)
		players = append(players, p)
	}
	return players, nil
}

/*
query the database for all seasons, populates global seasons store
example: seasonId: 22025 | season: 2024-25 | WSeason: 2025-26
*/
func UpdateSeasons(db *sql.DB) ([]Season, error) {
	e := errd.InitErr()
	rows, err := db.Query(pgdb.AllSeasons)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
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
	e := errd.InitErr()
	rows, err := db.Query(pgdb.Teams)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
	}

	var teams []Team
	for rows.Next() {
		var tm Team
		rows.Scan(&tm.League, &tm.TeamId, &tm.TeamAbbr, &tm.CityTeam)
		teams = append(teams, tm)
	}
	return teams, nil
}

// query db for team season records to populate records table
// moved to store rather than querying for every request
func UpdateTeamRecords(db *sql.DB, cs *CurrentSeasons) (TeamRecords, error) {
	e := errd.InitErr()
	var team_recs TeamRecords

	sl := cs.LgSznsByMonth(time.Now())
	rows, err := db.Query(pgdb.TeamSznRecords, sl.SznId, sl.WSznId)
	if err != nil {
		e.Msg = "error getting team records"
		return team_recs, e.BuildErr(err)
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
