package api

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
)

type TeamRecords struct {
	NBARecords  []TeamRecord `json:"nba_team_records"`
	WNBARecords []TeamRecord `json:"wnba_team_records"`
}

type TeamRecord struct {
	League     string `json:"league"`
	SeasonId   uint64 `json:"season_id"`
	Season     string `json:"season"`
	SeasonDesc string `json:"season_desc"`
	TeamId     uint64 `json:"team_id"`
	Team       string `json:"team"`
	TeamLong   string `json:"team_long"`
	Wins       uint16 `json:"wins"`
	Losses     uint16 `json:"losses"`
}

// query db for team season records to populate records table
func GetTeamRecords(db *sql.DB, cs *CurrentSeasons) (TeamRecords, error) {
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

// get from memory
func TeamRecordsJSON(team_recs *TeamRecords) ([]byte, error) {
	e := errd.InitErr()
	js, err := json.Marshal(team_recs)
	if err != nil {
		e.Msg = "error marshaling team records"
		return nil, e.BuildErr(err)
	}
	return js, nil
}
