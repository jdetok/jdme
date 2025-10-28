package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RemoveDiacritics(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	output, _, _ := transform.String(t, input)
	return output
}

// TODO: season map[uint64]map[uint64]string

// rewrite global stores of players, teams, etc using hash maps rather than arrays
type StMaps struct {
	PlayerIdDtl   map[uint64]*StPlayer // id as key, struct as val
	PlayerNameDtl map[string]*StPlayer // mame as key, struct as val
	PlayerIdName  map[uint64]string    // id as key, name as val
	PlayerNameId  map[string]uint64    // player name as key, id as val (for lookup)
	SeasonPlayers map[uint64]map[uint64]string
	SeasonPNames  map[int]map[string]struct{}
	Teams         map[uint64]*StTeam
	TeamIds       map[string]uint64
}

// player struct - to be stored as value in a map of player id keys
type StPlayer struct {
	Id      uint64
	Name    string
	Lowr    string // name all lower case
	Lg      string
	MaxRSzn uint64
	MinRSzn uint64
	MaxPSzn uint64
	MinPSzn uint64
	Teams   []uint64 // teams player has played for
}

func (sm *StMaps) MakeMaps() {
	sm.PlayerIdDtl = map[uint64]*StPlayer{}
	sm.PlayerNameDtl = map[string]*StPlayer{}
	sm.PlayerIdName = map[uint64]string{}
	sm.PlayerNameId = map[string]uint64{}
	sm.SeasonPlayers = map[uint64]map[uint64]string{}
	sm.SeasonPNames = map[int]map[string]struct{}{}
}

// team struct
// todo: hold old versions as well
type StTeam struct {
	Id   uint64
	Name string
	Lg   string
}

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapTeamIds(db *sql.DB) error {
	sm.TeamIds = map[string]uint64{}
	// get all team ids
	teams, err := db.Query("select distinct team_id from stats.tbox")
	if err != nil {
		return err
	}
	// convert each team id string to uint64
	for teams.Next() {
		var idstr string
		teams.Scan(&idstr)
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			return err
		}
		sm.TeamIds[idstr] = id
	}
	return nil
}

// query all player detail from db, map to id/name
func (sm *StMaps) MapPlayers(db *sql.DB) error {

	// get data for player struct (every player in lg.plr table)
	rows, err := db.Query(pgdb.QPlayerStore)
	if err != nil {
		return err
	}

	// each player creates a StPlayer struct that gets mapped to player id and
	// name also creates a map for each season and player (to determine whether a
	// player played in a particular season). also need to map teams
	//
	for rows.Next() {
		var p StPlayer
		var tms string  // comma separated string to be converted to []string
		var lowr string // run remove dia on each player's lowr
		rows.Scan(&p.Id, &p.Name, &lowr, &p.Lg, &p.MaxRSzn, &p.MinRSzn,
			&p.MaxPSzn, &p.MinPSzn, &tms)

		p.Lowr = RemoveDiacritics(lowr)
		// split tms string to slice of strings, get converted teamid uint from sm.TeamIds
		teamsStrArr := strings.SplitSeq(tms, ",")
		for t := range teamsStrArr {
			teamId := sm.TeamIds[t]
			p.Teams = append(p.Teams, teamId)
		}

		// iterate through each season played add the player to the season players map
		for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
			// init season players key for current season if it's nil
			if sm.SeasonPlayers[s] == nil {
				sm.SeasonPlayers[s] = map[uint64]string{}
			}
			// ADD PLAYER ID:NAME MAP TO SEASON MAP FOR EACH SEASON PLAYED
			sm.SeasonPlayers[s][p.Id] = p.Name

			if sm.SeasonPNames[int(s)] == nil {
				sm.SeasonPNames[int(s)] = map[string]struct{}{}
			}
			sm.SeasonPNames[int(s)][p.Lowr] = struct{}{}
		}

		// map player struct to id & name
		sm.PlayerIdDtl[p.Id] = &p
		sm.PlayerNameDtl[p.Lowr] = &p

		// map id to name & name to id
		sm.PlayerIdName[p.Id] = p.Lowr
		sm.PlayerNameId[p.Lowr] = p.Id
	}
	return nil
}

// check if player id is in season id map
func (sm *StMaps) PlayedInSzn(searchP, szn uint64) bool {
	_, ok := sm.SeasonPlayers[szn][searchP]
	if !ok {
		fmt.Printf("no match for season %d and player %d\n", szn, searchP)
		return false
	}
	fmt.Printf("%d played in %d\n", searchP, szn)
	return true
}

func (sm *StMaps) PlayerExists(searchP string) bool {
	_, ok := sm.PlayerNameId[searchP]
	return ok
}

func (sm *StMaps) PNameInSzn(searchP string, searchS int) bool {
	_, ok := sm.SeasonPNames[searchS][searchP]
	return ok
}
