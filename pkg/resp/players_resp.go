package resp

import (
	"encoding/json"
	"net/http"
)

type PlayerQuery struct {
	Player string
	Team   string
	Season string
	League string
}

type PlayerTeamSeason struct {
	Player string
	Team   string
	Season string
}
type PQueryIds struct {
	PId uint64
	TId uint64
	SId int
}

// outer response struct
type RespPlayerDash struct {
	Meta     RespMeta  `json:"request_meta"`
	Results  []RespObj `json:"player"`
	ErrorMsg string    `json:"error_string,omitempty"`
}

func NewRespPlayerDash(r *http.Request) *RespPlayerDash {
	return &RespPlayerDash{Meta: *NewRespMeta(r)}
}

func (rp *RespPlayerDash) WriteResp(w http.ResponseWriter) error {
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		return err
	}
	return nil
}

// each player's outermost struct, members of Resp.Results slice
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
// changed underscores in json tags to spaces
type RespPlayerStatsShtg struct {
	Fg  RespPlayerStatsShtgType `json:"field goals"`
	Fg3 RespPlayerStatsShtgType `json:"three pointers"`
	Ft  RespPlayerStatsShtgType `json:"free throws"`
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
