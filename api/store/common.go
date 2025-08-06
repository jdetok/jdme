package store

import (
	"database/sql"
	"strings"
	"unicode"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Player struct {
	PlayerId     uint64
	Name         string
	League       string
	SeasonIdMax  uint64
	SeasonIdMin  uint64
	PSeasonIdMax uint64
	PSeasonIdMin uint64
}

type Season struct {
	SeasonId string `json:"season_id"`
	Season   string `json:"season"`
	WSeason  string `json:"wseason"`
}

type Team struct {
	League   string `json:"league"`
	TeamId   string `json:"team_id"`
	TeamAbbr string `json:"team"`
	CityTeam string `json:"team_long"`
	LogoUrl  string `json:"-"`
}

// REMOVE NON SPACING CHARACTERS -- e.g. Dončić becomes doncic
func Unaccent(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	output, _, _ := transform.String(t, input)
	return output
}

// makes src url for team img
func (t Team) MakeLogoUrl() string {
	lg := strings.ToLower(t.League)
	return ("https://cdn." + lg + ".com/logos/" +
		lg + "/" + t.TeamId + "/primary/L/logo.svg")
}

// QUERY FOR PLAYER ID, PLAYER AND SAVE TO A LIST OF PLAYER STRUCTS
func GetPlayers(db *sql.DB) ([]Player, error) {
	e := errd.InitErr()
	// rows, err := db.Query(mdb.PlayersSeason.Q)
	rows, err := db.Query(pgdb.PlayersSeason.Q)
	if err != nil {
		e.Msg = "query failed"
		return nil, e.BuildErr(err)
	}
	var players []Player
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin, &p.PSeasonIdMax, &p.PSeasonIdMin)
		p.Name = Unaccent(p.Name) // REMOVE ACCENTS FROM NAMES
		players = append(players, p)
	}
	return players, nil
}

// seasons
func GetSeasons(db *sql.DB) ([]Season, error) {
	// fmt.Println("querying seasons & saving to struct")
	e := errd.InitErr()
	// rows, err := db.Query(mdb.RSeasons.Q)
	rows, err := db.Query(pgdb.AllSeasons.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
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
	e := errd.InitErr()
	// rows, err := db.Query(mdb.Teams.Q)
	rows, err := db.Query(pgdb.Teams.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
	}

	var teams []Team
	for rows.Next() {
		var tm Team
		rows.Scan(&tm.League, &tm.TeamId, &tm.TeamAbbr, &tm.CityTeam)
		teams = append(teams, tm)
	}
	return teams, nil
}
