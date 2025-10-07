package api

import (
	"encoding/json"

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
