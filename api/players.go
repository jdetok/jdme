package api

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

/*
primary database query function for the /players endpoint. queries the api
tables in the database sing the passed player, season, team ID to get the
player's stats. defaults to TeamTopScorerDash query, which gets the dash for
the top scorer of the most recent night's games. this is called when the site
loads. the response is scanned into the structs defined in resp.go, before being
marshalled into json and returned to write as the http response
*/
func (r *Resp) GetPlayerDash(db *sql.DB, pId uint64, sId uint64, tId uint64) ([]byte, error) {
	e := errd.InitErr()
	var q string
	var p uint64

	// if 0 is passed as tId, query by player_id. otherwise, query by team_id
	switch tId {
	case 0:
		logd.Logc(fmt.Sprintf("querying player_id: %d | season_id: %d", pId, sId))
		q = pgdb.PlayerDash.Q
		p = pId
	default:
		logd.Logc(fmt.Sprintf("querying team_id: %d | season_id: %d", tId, sId))
		q = pgdb.TeamTopScorerDash.Q
		p = tId
	}

	// QUERY SEASON PLAYERDASH FOR pId OR FOR TOP SCORER OF TEAM (tId) PASSED
	rows, err := db.Query(q, p, sId)
	if err != nil {
		e.Msg = "error during player dash query"
		return nil, e.BuildErr(err)
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
		rp.HndlRespRow(&p, &s)
	}
	// handle aggregate season ids (all, regular season, playoffs)
	HndlAggsIds(&rp.Meta.SeasonId, &rp.Meta.StatType)

	// assign nba or wnba season only based on league
	t.HndlSeason(&rp.Meta.League, &rp.Meta.Season)

	// build table captions & image urls
	rp.Meta.MakeCaptions()
	rp.Meta.MakeHeadshotUrl()
	rp.Meta.MakeTeamLogoUrl()
	r.Results = append(r.Results, rp)

	// marshal response & return json []byte
	js, err := json.Marshal(r)
	if err != nil {
		e.Msg = "failed to marshal structs to json"
		return nil, e.BuildErr(err)
	}
	return js, nil
}

/*
switch between totals (sums) and pergame (averages) stats based on the
Meta.StatType field
*/
func (rp *RespObj) HndlRespRow(p *RespPlayerSznOvw, s *RespPlayerStats) {
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

// accept pointers of league and season, switch season/wseason on league
func (t *RespSeasonTmp) HndlSeason(league *string, season *string) {
	switch *league {
	case "NBA":
		*season = t.Season
	case "WNBA":
		*season = t.WSeason
	}
}

/*
accept pointers of season_id and stat type, switch season to handle stat type
used for aggregate seasons, deprecated
*/
func HndlAggsIds(sId *uint64, sType *string) {
	switch *sId {
	case 99999:
		*sType = "career regular season + playoffs"
	case 99998:
		*sType = "career regular season"
	case 99997:
		*sType = "career playoffs"
	default:
		*sType = "regular season"
	}
}

func (m *RespPlayerMeta) MakeCaptions() {
	m.Caption = fmt.Sprintf("%s - %s", m.Player, m.TeamName)
	m.CaptionShort = fmt.Sprintf("%s - %s", m.Player, m.Team)
	m.BoxCapTot = fmt.Sprintf("Box Totals - %s", m.Season)
	m.BoxCapAvg = fmt.Sprintf("Box Averages - %s", m.Season)
	m.ShtgCapTot = fmt.Sprintf("Shooting Totals - %s", m.Season)
	m.ShtgCapAvg = fmt.Sprintf("Shooting Averages - %s", m.Season)
}
