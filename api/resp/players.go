package resp

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/applog"
	"github.com/jdetok/go-api-jdeko.me/mdb"
)

func (r *Resp) GetPlayerDash(db *sql.DB, pId uint64, sId uint64, tId uint64) ([]byte, error) {
	e := applog.AppErr{Process: "GetPlayerDash()"}
	var q string
	var p uint64
	switch tId {
	case 0:
		q = mdb.Player.Q
		p = pId
	default:
		q = mdb.TeamSeasonTopP.Q
		p = tId
	}

	rows, err := db.Query(q, p, sId)
	if err != nil {
		e.Msg = fmt.Sprintf(
			`player dash query (player_id: %d | season_id: %d)`, pId, sId)
		return nil, e.BuildError(err)
	}
	var t RespSeasonTmp // temp seasons for NBA/WNBA, handled after loop
	var rp RespObj
	for rows.Next() {
		// temp structs, handled in hndlRespRow
		var s RespPlayerStats
		var p RespPlayerSznOvw
		rows.Scan( // MUST BE IN ORDER OF QUERY
			&rp.Meta.PlayerId, &rp.Meta.TeamId, &rp.Meta.League,
			&rp.Meta.SeasonId, &rp.Meta.StatType, &rp.Meta.Player,
			&rp.Meta.Team, &rp.Meta.TeamName,
			&rp.SeasonOvw.GamesPlayed, &p.Minutes,
			&s.Box.Points, &s.Box.Assists, &s.Box.Rebounds,
			&s.Box.Steals, &s.Box.Blocks,
			&s.Shtg.Fg.Makes, &s.Shtg.Fg.Attempts, &s.Shtg.Fg.Percent,
			&s.Shtg.Fg3.Makes, &s.Shtg.Fg3.Attempts, &s.Shtg.Fg3.Percent,
			&s.Shtg.Ft.Makes, &s.Shtg.Ft.Attempts, &s.Shtg.Ft.Percent,
			&t.Season, &t.WSeason)
		// switch on stat type to assign stats to appropriate struct
		rp.hndlRespRow(&p, &s)
	}
	// handle aggregate season ids (all, regular season, playoffs)
	hndlAggsIds(&rp.Meta.SeasonId, &rp.Meta.StatType)
	t.hndlSeason(&rp.Meta.League, &rp.Meta.Season)

	// build table captions & image urls
	rp.Meta.MakeCaptions()
	rp.Meta.MakeHeadshotUrl()
	rp.Meta.MakeTeamLogoUrl()
	r.Results = append(r.Results, rp)

	js, err := json.Marshal(r)
	if err != nil {
		e.Msg = "failed to marshal structs to json"
		return nil, e.BuildError(err)
	}
	return js, nil
}

func (rp *RespObj) hndlRespRow(p *RespPlayerSznOvw, s *RespPlayerStats) {
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
func (t *RespSeasonTmp) hndlSeason(league *string, season *string) {
	switch *league {
	case "NBA":
		*season = t.Season
	case "WNBA":
		*season = t.WSeason
	}
}

// accept pointers of season_id and stat type, switch season to handle stat type
func hndlAggsIds(sId *uint64, sType *string) {
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
