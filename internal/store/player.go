package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/internal/errs"
	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

/* INTENT:
create player structs that will init when a player is searched
should be able to handle searches for the following:
  - player only
  - player/season
  - player/team
  - player/season/team

the shooting stats are broke out into three separate struct with made,
atp, pct for each type of shot. these wrap into another struct
*/

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
	StatType     string `json:"stat_type"`
	Player       string `json:"player"`
	Team         string `json:"team"`
	TeamName     string `json:"team_name"`
	Caption      string `json:"caption"`
	CaptionShort string `json:"caption_short"`
	HeadshotUrl  string `json:"headshot_url"`
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
	Fg  RespPlayerStatsFg  `json:"fg"`
	Fg3 RespPlayerStatsFg3 `json:"fg3"`
	Ft  RespPlayerStatsFt  `json:"ft"`
}

// change these to just made atpt pct cause putting in parent struct
type RespPlayerStatsFg struct {
	Makes    float32 `json:"m"`
	Attempts float32 `json:"a"`
	Percent  string  `json:"p"`
}
type RespPlayerStatsFg3 struct {
	Makes    float32 `json:"m"`
	Attempts float32 `json:"a"`
	Percent  string  `json:"p"`
}
type RespPlayerStatsFt struct {
	Makes    float32 `json:"m"`
	Attempts float32 `json:"a"`
	Percent  string  `json:"p"`
}

// move caps and url to different struct

func (m *RespPlayerMeta) MakeCaptions() {
	m.Caption = fmt.Sprintf("%s - %s", m.Player, m.TeamName)
	m.CaptionShort = fmt.Sprintf("%s - %s", m.Player, m.Team)
}

func (m *RespPlayerMeta) MakeHeadshotUrl() {
	lg := strings.ToLower(m.League)
	pId := strconv.Itoa(int(m.PlayerId))
	m.HeadshotUrl = fmt.Sprintf(
		`https://cdn.%s.com/headshots/%s/latest/1040x760/%s.png`,
		lg, lg, pId)
}

// DB QUERY
func (r *Resp) GetPlayerDash(db *sql.DB, pId uint64, sId uint64) ([]byte, error) {
	e := errs.ErrInfo{Prefix: "getting player dash"}
	rows, err := db.Query(mariadb.Player.Q, pId, sId)
	if err != nil {
		e.Msg = fmt.Sprintf(
			`player dash query (player_id: %d | season_id: %d)`, pId, sId)
		return nil, e.Error(err)
	}

	var rp RespObj
	for rows.Next() {
		// temp struct to add logic to which stat struct is populated
		var s RespPlayerStats
		var p RespPlayerSznOvw
		rows.Scan(
			&rp.Meta.PlayerId, &rp.Meta.TeamId, &rp.Meta.League,
			&rp.Meta.SeasonId, &rp.Meta.StatType, &rp.Meta.Player,
			&rp.Meta.Team, &rp.Meta.TeamName,
			&rp.SeasonOvw.GamesPlayed, &p.Minutes,
			&s.Box.Points, &s.Box.Assists, &s.Box.Rebounds,
			&s.Box.Steals, &s.Box.Blocks,
			&s.Shtg.Fg.Makes, &s.Shtg.Fg.Attempts, &s.Shtg.Fg.Percent,
			&s.Shtg.Fg3.Makes, &s.Shtg.Fg3.Attempts, &s.Shtg.Fg3.Percent,
			&s.Shtg.Ft.Makes, &s.Shtg.Ft.Attempts, &s.Shtg.Ft.Percent)
		rp.Meta.MakeCaptions()
		rp.Meta.MakeHeadshotUrl()

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

	switch rp.Meta.SeasonId {
	case 99999:
		rp.Meta.StatType = "career regular season + playoffs"
	case 99998:
		rp.Meta.StatType = "career regular season"
	case 99997:
		rp.Meta.StatType = "career playoffs"
	default:
		rp.Meta.StatType = "regular season"
	}

	r.Results = append(r.Results, rp)

	js, err := json.Marshal(r)
	if err != nil {
		e.Msg = "failed to marshal structs to json"
	}

	return js, nil
}
