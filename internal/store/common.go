package store

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/internal/errs"
	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

type GameMeta struct {
	SeasonID string `json:"season_id"`
	GameId   string `json:"game_id"`
	GameDate string `json:"game_date"`
}

type PlayerMeta struct {
	PlayerId     string `json:"player_id"`
	TeamId       string `json:"team_id"`
	League       string `json:"league"`
	Player       string `json:"player"`
	Team         string `json:"team"`
	TeamName     string `json:"team_name"`
	Caption      string `json:"caption"`
	CaptionShort string `json:"caption_short"`
	HeadshotUrl  string `json:"headshot_url"`
}

// idea: break out box and shooting
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

type BoxStats struct {
	Minutes  string `json:"minutes"`
	Points   string `json:"points"`
	Assists  string `json:"assists"`
	Rebounds string `json:"rebounds"`
	Steals   string `json:"steals"`
	Blocks   string `json:"blocks"`
}

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

type Player struct {
	PlayerId    uint64
	Name        string
	League      string
	SeasonIdMax uint64
	SeasonIdMin uint64
}

type Season struct {
	SeasonId string
	Season   string
	WSeason  string
}

type Team struct {
	League   string
	TeamId   string
	TeamAbbr string
	CityTeam string
	LogoUrl  string
}

func (pm *PlayerMeta) MakeCaptions() {
	pm.Caption = fmt.Sprintf("%s - %s", pm.Player, pm.TeamName)
	pm.CaptionShort = fmt.Sprintf("%s - %s", pm.Player, pm.Team)
}

func (pm *PlayerMeta) MakeHeadshotUrl() {
	lg := strings.ToLower(pm.League)
	pm.HeadshotUrl = fmt.Sprintf(
		`https://cdn.%s.com/headshots/%s/latest/1040x760/%s.png`,
		lg, lg, pm.PlayerId)
}

// makes src url for team img
func (t Team) MakeLogoUrl() string {
	lg := strings.ToLower(t.League)
	return ("https://cdn." + lg + ".com/logos/" +
		lg + "/" + t.TeamId + "/primary/L/logo.svg")
}

// QUERY FOR PLAYER ID, PLAYER AND SAVE TO A LIST OF PLAYER STRUCTS
func GetPlayers(db *sql.DB) ([]Player, error) {
	fmt.Println("querying players & saving to struct")
	e := errs.ErrInfo{Prefix: "saving players to structs"}
	rows, err := db.Query(mariadb.PlayersSeason.Q)
	if err != nil {
		e.Msg = "query failed"
		return nil, e.Error(err)
	}
	var players []Player
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin)
		// convert to lowercase to match requests
		p.Name = strings.ToLower(p.Name)
		p.League = strings.ToLower(p.League)
		players = append(players, p)
	}
	return players, nil
}
func randPlayer(players []Player) uint64 {
	numPlayers := len(players)
	randNum := rand.IntN(numPlayers)
	return players[randNum].PlayerId
}
func GetpIdsId(players []Player, player string, seasonId string) (uint64, uint64) {
	sId, _ := strconv.ParseUint(seasonId, 10, 32)
	var pId uint64

	if player == "random" {
		pId = randPlayer(players)
	} else if _, err := strconv.ParseUint(player, 10, 64); err == nil {
		pId, _ = strconv.ParseUint(player, 10, 64)
	} else {
		for _, p := range players {
			if p.Name == player { // return match playerid (uint32) as string
				pId = p.PlayerId
			}
		}
	}

	for _, p := range players {
		if p.PlayerId == pId { // return match playerid (uint32) as string
			if sId > p.SeasonIdMax {
				return pId, p.SeasonIdMax
			} else if sId < p.SeasonIdMin {
				return pId, p.SeasonIdMin
			} else {
				return pId, sId
			}
		}
	}

	return pId, sId
}

func SearchPlayers(players []Player, pSearch string) string {
	for _, p := range players {
		if p.Name == pSearch { // return match playerid (uint32) as string
			return strconv.FormatUint(p.PlayerId, 10)
		}
	}
	return ""
}

// seasons
func GetSeasons(db *sql.DB) ([]Season, error) {
	fmt.Println("querying seasons & saving to struct")
	e := errs.ErrInfo{Prefix: "saving seasons to struct"}
	rows, err := db.Query(mariadb.Seasons.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.Error(err)
	}

	var seasons []Season
	for rows.Next() {
		var szn Season
		rows.Scan(&szn.SeasonId, &szn.Season, &szn.WSeason)
		seasons = append(seasons, szn)
	}

	return seasons, nil
}

// teams
func GetTeams(db *sql.DB) ([]Team, error) {
	fmt.Println("querying teams & saving to struct")
	e := errs.ErrInfo{Prefix: "saving teams to struct"}
	rows, err := db.Query(mariadb.Teams.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.Error(err)
	}

	var teams []Team
	for rows.Next() {
		var tm Team
		rows.Scan(&tm.League, &tm.TeamId, &tm.TeamAbbr, &tm.CityTeam)
		teams = append(teams, tm)
	}
	return teams, nil
}
