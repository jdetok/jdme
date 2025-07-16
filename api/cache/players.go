package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/applog"
	"github.com/jdetok/go-api-jdeko.me/mdb"
)

// outermost struct, returned to http handler as json string
type Resp struct {
	Results []RespObj `json:"player"`
}

type RespObj struct {
	Meta      RespPlayerMeta   `json:"player_meta"`
	SeasonOvw RespPlayerSznOvw `json:"playtime"`
	Totals    RespPlayerStats  `json:"totals"`
	PerGame   RespPlayerStats  `json:"per_game"`
}

type RespPlayerMeta struct {
	PlayerId     uint64 `json:"player_id"`
	TeamId       uint64 `json:"team_id"`
	League       string `json:"league"`
	SeasonId     uint64 `json:"season_id"`
	StatType     string `json:"-"`
	Player       string `json:"player"`
	Team         string `json:"team"`
	TeamName     string `json:"team_name"`
	Season       string `json:"season"`
	Caption      string `json:"caption"`
	CaptionShort string `json:"caption_short"`
	BoxCapTot    string `json:"cap_box_tot"`
	BoxCapAvg    string `json:"cap_box_avg"`
	ShtgCapTot   string `json:"cap_shtg_tot"`
	ShtgCapAvg   string `json:"cap_shtg_avg"`
	HeadshotUrl  string `json:"headshot_url"`
	TeamLogoUrl  string `json:"team_logo_url"`
}

type RespPlayerSznOvw struct {
	GamesPlayed   uint16  `json:"games_played"`
	Minutes       float32 `json:"minutes"`
	MinutsPerGame float32 `json:"minutes_pg"`
}

type RespPlayerStats struct {
	Box  RespPlayerStatsBox  `json:"box_stats"`
	Shtg RespPlayerStatsShtg `json:"shooting"`
}

type RespPlayerStatsBox struct {
	Points   float32 `json:"points"`
	Assists  float32 `json:"assists"`
	Rebounds float32 `json:"rebounds"`
	Steals   float32 `json:"steals"`
	Blocks   float32 `json:"blocks"`
}

// struct to wrap shooting stats
type RespPlayerStatsShtg struct {
	Fg  RespPlayerStatsShtgType `json:"field_goals"`
	Fg3 RespPlayerStatsShtgType `json:"three_pointers"`
	Ft  RespPlayerStatsShtgType `json:"free_throws"`
}

// change these to just made atpt pct cause putting in parent struct
type RespPlayerStatsShtgType struct {
	Makes    float32 `json:"made"`
	Attempts float32 `json:"attempted"`
	Percent  string  `json:"percentage"`
}

// temporary struct used in GetPlayerDash
type RespSeasonTmp struct {
	Season  string
	WSeason string
}

/*
CREATE JSON RESPONSE FOR /player
accept player/team/season ids & query api table in database
scan rows to structs, build table captions & image urls,
marshal & return structured json string
*/
// func (r *Resp) GetPlayerDash(db *sql.DB, pId uint64, sId uint64, tId uint64) ([]byte, error) {
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
	// rows, err := db.Query(mariadb.Player.Q, pId, sId)
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

func (m *RespPlayerMeta) MakeHeadshotUrl() {
	lg := strings.ToLower(m.League)
	pId := strconv.Itoa(int(m.PlayerId))
	m.HeadshotUrl = fmt.Sprintf(
		`https://cdn.%s.com/headshots/%s/latest/1040x760/%s.png`,
		lg, lg, pId)
}

func (m *RespPlayerMeta) MakeTeamLogoUrl() {
	lg := strings.ToLower(m.League)
	tId := strconv.Itoa(int(m.TeamId))
	m.TeamLogoUrl = fmt.Sprintf(
		`https://cdn.%s.com/logos/%s/%s/primary/L/logo.svg`,
		lg, lg, tId)
}
