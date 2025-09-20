package api

import (
	"fmt"
	"strconv"
	"strings"
)

// outer response struct
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

// all player statistics (strings)
type Stats struct {
	Minutes  string `json:"minutes"`
	Points   string `json:"points"`
	Assists  string `json:"assists"`
	Rebounds string `json:"rebounds"`
	Steals   string `json:"steals"`
	Blocks   string `json:"blocks"`
	FgMade   string `json:"fg_made"`
	FgAtpt   string `json:"fg_atpt"`
	FgPct    string `json:"fg_pct"`
	Fg3Made  string `json:"fg3_made"`
	Fg3Atpt  string `json:"fg3_atpt"`
	Fg3Pct   string `json:"fg3_pct"`
	FtMade   string `json:"ft_made"`
	FtAtpt   string `json:"ft_atpt"`
	FtPct    string `json:"ft_pct"`
}

// only box stats (strings)
type BoxStats struct {
	Minutes  string `json:"minutes"`
	Points   string `json:"points"`
	Assists  string `json:"assists"`
	Rebounds string `json:"rebounds"`
	Steals   string `json:"steals"`
	Blocks   string `json:"blocks"`
}

// only shooting stats (strings)
type ShootingStats struct {
	FgMade  string `json:"fg_made"`
	FgAtpt  string `json:"fg_atpt"`
	FgPct   string `json:"fg_pct"`
	Fg3Made string `json:"fg3_made"`
	Fg3Atpt string `json:"fg3_atpt"`
	Fg3Pct  string `json:"fg3_pct"`
	FtMade  string `json:"ft_made"`
	FtAtpt  string `json:"ft_atpt"`
	FtPct   string `json:"ft_pct"`
}

// outermost struct, returned to http handler as json string
type RespPlayerSznOvw struct {
	GamesPlayed   uint16  `json:"games_played"`
	Minutes       float32 `json:"minutes"`
	MinutsPerGame float32 `json:"minutes_pg"`
}

// struct that holds both box and shooting stats
type RespPlayerStats struct {
	Box  RespPlayerStatsBox  `json:"box_stats"`
	Shtg RespPlayerStatsShtg `json:"shooting"`
}

// box stats as floats
type RespPlayerStatsBox struct {
	Points   float32 `json:"points"`
	Assists  float32 `json:"assists"`
	Rebounds float32 `json:"rebounds"`
	Steals   float32 `json:"steals"`
	Blocks   float32 `json:"blocks"`
}

// shooting stats wrapper struct - holds shtg type structs for twos, three, free throws
type RespPlayerStatsShtg struct {
	Fg  RespPlayerStatsShtgType `json:"field_goals"`
	Fg3 RespPlayerStatsShtgType `json:"three_pointers"`
	Ft  RespPlayerStatsShtgType `json:"free_throws"`
}

// struct to hold a category of shooting stats - should be stored in wrapper struct
// for twos, threes, free throws
type RespPlayerStatsShtgType struct {
	Makes    float32 `json:"made"`
	Attempts float32 `json:"attempted"`
	Percent  string  `json:"percentage"`
}

// temporary struct used in GetPlayerDash to assign appropriate league for each player
type RespSeasonTmp struct {
	Season  string
	WSeason string
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
func (m *RespPlayerMeta) MakeTeamLogoUrl() {
	lg := strings.ToLower(m.League)
	tId := strconv.Itoa(int(m.TeamId))
	m.TeamLogoUrl = fmt.Sprintf(
		`https://cdn.%s.com/logos/%s/%s/primary/L/logo.svg`,
		lg, lg, tId)
}
