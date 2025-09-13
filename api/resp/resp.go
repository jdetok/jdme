package resp

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/store"
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

// accept the slice of all players and a seasonId, return a slice with just the
// active players from the passed season id
func slicePlayersSzn(players []store.Player, seasonId uint64) ([]store.Player, error) {
	var plslice []store.Player
	for _, p := range players { // EXPAND THIS IF TO CATCH PLAYOFF SEASONS AS WELL

		// if nba & pre october, subtract 1 from season id
		var sId uint64 = seasonId
		if uint64(time.Now().Month()) <= 10 {
			if p.League == "nba" {
				sId--
			}
		}

		// append players to the random slice if the passed season id between player min and max season
		if (sId >= 20000 && sId < 30000) && (sId <= p.SeasonIdMax && sId >= p.SeasonIdMin) || (sId >= 40000 && sId < 50000) && (sId <= p.PSeasonIdMax && sId >= p.PSeasonIdMin) {
			plslice = append(plslice, p)
		} else if sId >= 88888 {
			plslice = append(plslice, p)
		}
	}
	return plslice, nil
}

/*
accept slice of Player structs and a season id, call slicePlayerSzn to create
a new slice with only players from the specified season. then, generate a
random number and return the player at that index in the slice
*/
func randPlayer(pl []store.Player, sId uint64) uint64 {
	players, _ := slicePlayersSzn(pl, sId)
	numPlayers := len(players)
	randNum := rand.IntN(numPlayers)
	return players[randNum].PlayerId
}

/*
player name and season ID from get request passed here, returns the player's
ID and the season ID. if 'player' variable == "random", the randPlayer function
is called. a player ID also can be passed as the player parameter, it will just
be converted to an int and returned
*/
func GetpIdsId(players []store.Player, player string, seasonId string) (uint64, uint64) {
	sId, _ := strconv.ParseUint(seasonId, 10, 32)
	var pId uint64

	if player == "random" { // call randplayer function
		pId = randPlayer(players, sId)
	} else if _, err := strconv.ParseUint(player, 10, 64); err == nil {
		// if it's numeric keep it and convert to uint64
		pId, _ = strconv.ParseUint(player, 10, 64)
	} else { // search name through players list
		for _, p := range players {
			if p.Name == player { // return match playerid (uint32) as string
				pId = p.PlayerId
			}
		}
	}

	// loop through players to check that queried season is within min-max seasons
	for _, p := range players {
		if p.PlayerId == pId {
			return pId, HandlesId(sId, &p)
		}
	}
	return pId, sId
}

/*
accept a season id and a pointer to a Player struct, validate the player was active
in the passed season, return a valid season ID if not. if season id starts with an
8 the player's max regular season will be returned. if it starts with a 7, their
max playoff season will be returned. if it starts with a 4, it will first verify
player has played in a playoff game, and will return their max regular season if
they haven't. a season id starting with 2 will return a regular season. for both
regular season and playoffs, the function will verify the player played in said
season, and return either their max or min (whichever is closer) season  if they
did not
*/
func HandlesId(sId uint64, p *store.Player) uint64 {
	if strconv.FormatUint(sId, 10)[1:] == "9999" { // agg seasons
		return sId
	} else if sId >= 80000 && sId < 90000 {
		return p.SeasonIdMax // return most recent season
	} else if sId >= 70000 && sId < 80000 {
		return p.PSeasonIdMax // return most recent season
	} else if sId >= 40000 && sId < 50000 {
		if p.PSeasonIdMax < 40000 { // player has no playeroff, return max reg season
			return p.SeasonIdMax // return reg season if player has no playoffs
		}
		if sId == 49999 { // playoff career
			return sId
		}

		if sId > p.PSeasonIdMax {
			return p.PSeasonIdMax
		}
		if sId < p.PSeasonIdMin {
			return p.PSeasonIdMin
		}
	} else if sId >= 20000 && sId < 30000 {
		if sId > p.SeasonIdMax {
			if sId == 29999 { // reg season career
				return sId
			}
			return p.SeasonIdMax
		}
		if sId < p.SeasonIdMin {
			return p.SeasonIdMin
		}
	}
	return sId
}

// search player by name, return player id int if found
func SearchPlayers(players []store.Player, pSearch string) string {
	for _, p := range players {
		if p.Name == pSearch { // return match playerid (uint32) as string
			return strconv.FormatUint(p.PlayerId, 10)
		}
	}
	return ""
}
