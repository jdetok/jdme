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
	TeamIds       map[string]uint64
	// new ones
	PlrIds       map[uint64]struct{}
	PlrNms       map[string]struct{}
	SeasonPlrNms map[int]map[string]struct{}
	SeasonPlrIds map[int]map[uint64]string
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

func (sm *StMaps) MakeMaps(db *sql.DB) {
	sm.PlayerIdDtl = map[uint64]*StPlayer{}
	sm.PlayerNameDtl = map[string]*StPlayer{}
	sm.PlayerIdName = map[uint64]string{}
	sm.PlayerNameId = map[string]uint64{}
	sm.SeasonPlayers = map[uint64]map[uint64]string{}
	sm.SeasonPlrNms = map[int]map[string]struct{}{}
	sm.PlrIds = map[uint64]struct{}{}
	sm.PlrNms = map[string]struct{}{}
	if err := sm.MapTeamIds(db); err != nil {
		fmt.Println(err)
	}
	if err := sm.MapSeasons(db); err != nil {
		fmt.Println(err)
	}
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

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapSeasons(db *sql.DB) error {
	sm.SeasonPlrNms = map[int]map[string]struct{}{}
	sm.SeasonPlrIds = map[int]map[uint64]string{}

	// to handle season id = 0
	sm.SeasonPlrNms[0] = map[string]struct{}{}
	sm.SeasonPlrIds[0] = map[uint64]string{}

	// get all season ids
	szns, err := db.Query("select distinct szn_id from stats.tbox")
	if err != nil {
		return err
	}
	// convert each season id string to int
	for szns.Next() {
		var sznStr string
		szns.Scan(&sznStr)
		szn, err := strconv.Atoi(sznStr)
		if err != nil {
			return err
		}
		// create empty map for each season
		sm.SeasonPlrIds[szn] = map[uint64]string{}
		sm.SeasonPlrNms[szn] = map[string]struct{}{}
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
	for rows.Next() {
		var p StPlayer
		sm.MapPlayerRow(rows, &p)
	}
	return nil
}

// scan player results from db to maps
func (sm *StMaps) MapPlayerRow(rows *sql.Rows, p *StPlayer) {
	var tms string  // comma separated string to be converted to []string
	var lowr string // run remove dia on each player's lowr
	rows.Scan(&p.Id, &p.Name, &lowr, &p.Lg, &p.MaxRSzn, &p.MinRSzn,
		&p.MaxPSzn, &p.MinPSzn, &tms)

	// remove accents from lower case name
	p.Lowr = RemoveDiacritics(lowr)

	// split tms string from db to slice of strings
	sm.CleanTeamsSlice(p, tms)

	// iterate through each season played add the player to the season players map
	sm.MapPlrToSzn(p)

	// basic player id and name maps for simple Exists funcs
	sm.MapPlrIdsNms(p)

	// map player struct to id & name
	sm.PlayerIdDtl[p.Id] = p
	sm.PlayerNameDtl[p.Lowr] = p

	// map id to name & name to id
	sm.PlayerIdName[p.Id] = p.Lowr
	sm.PlayerNameId[p.Lowr] = p.Id
}

// insert player id and cleaned player name as keys in PlrIds and PlrNms maps
func (sm *StMaps) MapPlrIdsNms(p *StPlayer) {
	sm.PlrIds[p.Id] = struct{}{}
	sm.PlrNms[p.Lowr] = struct{}{}
}

// split comma separated string from db into slice of strings
func (sm *StMaps) CleanTeamsSlice(p *StPlayer, tms string) {
	// split tms string from db to slice of strings
	teamsStrArr := strings.SplitSeq(tms, ",")

	// iterate through each team player has played for
	// TODO: map players to team map similar to season maps
	for t := range teamsStrArr {
		// use TeamIds map created with MakeMaps() get the uint64 version of t
		// append to teams slice
		teamId := sm.TeamIds[t]
		p.Teams = append(p.Teams, teamId)
	}
}

// iterate through each season played add the player to the season players map
func (sm *StMaps) MapPlrToSzn(p *StPlayer) {
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrNms[int(s)][p.Lowr] = struct{}{}
		sm.SeasonPlrIds[int(s)][p.Id] = p.Lowr
	}
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

// search SeasonPlrNms or SeasonPlrIds
func (sm *StMaps) PNameInSzn(searchP string, searchS int) bool {
	var ok bool
	plrIdInt, err := strconv.ParseUint(searchP, 10, 64)
	if err != nil { // lookup by string if it fails to convert to int
		if searchS == 0 { // only verify player name if 0 season
			return sm.PlrNmExists(searchP)
		}
		_, ok = sm.SeasonPlrNms[searchS][searchP]
	} else { // lookup by player id if searchP is numeric (successful ParseUint)
		if searchS == 0 { // only verify player id if 0 season
			return sm.PlrIdExists(plrIdInt)
		}
		_, ok = sm.SeasonPlrIds[searchS][plrIdInt]
	}
	return ok
}

// check if passed player id exists in PlrIds map
func (sm *StMaps) PlrIdExists(searchP uint64) bool {
	_, ok := sm.PlrIds[searchP]
	return ok
}

// check if passed player name exists in PlrNms map
func (sm *StMaps) PlrNmExists(searchP string) bool {
	_, ok := sm.PlrNms[searchP]
	return ok
}
