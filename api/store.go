package api

import (
	"database/sql"
	"sync"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
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
func (app *App) CheckInMemStructs(interval, threshold time.Duration) {
	// call update func on intial run
	if app.Started == 0 {
		app.UpdateStructs()

		app.Started = 1
	}

	// update structs every interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(app.LastUpdate) > threshold {
			app.Lg.Infof("refreshing in-mem store")

			var wg sync.WaitGroup

			// calculate current NBA/WNBA seasons
			app.Store.CurrentSzns.GetCurrentSzns(time.Now())

			// update in memory structs
			wg.Add(1)
			go func(wg *sync.WaitGroup, app *App) {
				defer wg.Done()
				app.UpdateStructs()
			}(&wg, app)

			// update maps
			wg.Add(1)
			go func(wg *sync.WaitGroup, app *App) {
				defer wg.Done()
				if err := app.MStore.Rebuild(app.Database, app.Lg); err != nil {
					app.Lg.Errorf("failed to update player map")
				}
			}(&wg, app)
			wg.Wait()
		}
	}
}

// update players, seasons, and teams in memory structs slices
func (app *App) UpdateStructs() {
	app.Lg.Infof("updating in memory structs")

	var errP error
	msgP := "updating players structs"
	app.Store.Players, errP = UpdatePlayers(app.Database)
	if errP != nil {
		app.Lg.Errorf("failed %s\n%v", msgP, errP)
	}

	// update in memory seasons slice
	var errS error
	msgS := "updating seasons structs"
	app.Store.Seasons, errS = UpdateSeasons(app.Database)
	if errS != nil {
		app.Lg.Errorf("failed %s\n%v", msgS, errS)
	}

	// update in memory teams slice
	var errT error
	msgE := "updating teams structs"
	app.Store.Teams, errT = UpdateTeams(app.Database)
	if errT != nil {
		app.Lg.Errorf("failed %s\n%v", msgE, errP)
	}
	// update team records
	var errTR error
	msgTR := "updating team records"
	app.Store.TeamRecs, errTR = UpdateTeamRecords(app.Database, &app.Store.CurrentSzns)
	if errTR != nil {
		app.Lg.Errorf("failed %s\n%v", msgTR, errTR)
	}
	// update league top players
	var errLP error
	msgLP := "updating league top players struct"
	app.Store.TopLgPlayers, errLP = QueryTopLgPlayers(app.Database, &app.Store.CurrentSzns, "50")
	if errLP != nil {
		app.Lg.Errorf("failed %s\n%v", msgLP, errLP)
	}

	// update last update time
	updateTime := time.Now()
	app.LastUpdate = updateTime
	app.Lg.Infof("finished refreshing store")
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

// query db for team season records to populate records table
// moved to store rather than querying for every request
func UpdateTeamRecords(db *sql.DB, cs *CurrentSeasons) (TeamRecords, error) {
	var team_recs TeamRecords

	sl := cs.LgSznsByMonth(time.Now())
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
