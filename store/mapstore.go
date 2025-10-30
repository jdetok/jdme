package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
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

type MapStore struct {
	Maps *StMaps
	mu   sync.RWMutex
}

func (ms *MapStore) Set(newMaps *StMaps) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Maps = newMaps
}

func (ms *MapStore) Rebuild(db *sql.DB) error {
	fmt.Println("Rebuilding StMaps...")
	temp := MakeMaps(db)

	if err := temp.MapTeamIds(db); err != nil {
		return err
	}
	if err := temp.MapSeasons(db); err != nil {
		return err
	}
	if err := temp.MapPlayers(db); err != nil {
		fmt.Println("Error in MapPlayers:", err)
		return err
	}

	fmt.Printf("Rebuild complete: %d players mapped\n", len(temp.PlayerIdDtl))
	ms.Set(temp)
	return nil
}

func (ms *MapStore) Setup(db *sql.DB) {
	fmt.Println("MAP STORE SETUP STARTED")
	ms.Set(MakeMaps(db))
}

// rewrite global stores of players, teams, etc using hash maps rather than arrays
type StMaps struct {
	PlayerIdDtl   map[uint64]*StPlayer // id as key, struct as val
	PlayerNameDtl map[string]*StPlayer // mame as key, struct as val
	PlayerIdName  map[uint64]string    // id as key, name as val
	PlayerNameId  map[string]uint64    // player name as key, id as val (for lookup)
	SeasonPlayers map[int]map[uint64]string
	TeamIds       map[string]uint64
	// new ones
	PlrIds       map[uint64]struct{}
	PlrNms       map[string]struct{}
	SeasonPlrNms map[int]map[string]uint64
	SeasonPlrIds map[int]map[uint64]string

	// season team player map
	// example: SznTmPlrIds[szn][tm][plr]
	SznTmPlrIds map[int]map[uint64]map[uint64]string
}

// player struct - to be stored as value in a map of player id keys
type StPlayer struct {
	Id          uint64
	Name        string
	Lowr        string // name all lower case
	Lg          string
	MaxRSzn     int
	MinRSzn     int
	MaxPSzn     int
	MinPSzn     int
	Teams       []uint64 // teams player has played for
	TmsRostered map[uint64]uint64
}

// INITIAL MAP SETUP: must create empty maps before attempting to insert keys
// calls MapTeams and MapSeasons to setup an empty map for
// each season and team nested map
func MakeMaps(db *sql.DB) *StMaps {
	fmt.Println("creating empty maps")
	var sm StMaps
	sm.PlayerIdDtl = map[uint64]*StPlayer{}
	sm.PlayerNameDtl = map[string]*StPlayer{}
	sm.PlayerIdName = map[uint64]string{}
	sm.PlayerNameId = map[string]uint64{}
	sm.SeasonPlayers = map[int]map[uint64]string{}

	// map of seasons with nested map of player ids/names (cleaned)
	sm.SeasonPlrNms = map[int]map[string]uint64{}
	sm.SeasonPlrIds = map[int]map[uint64]string{}

	// map of player ids & player names to verify ONLY if player exists in the db
	sm.PlrIds = map[uint64]struct{}{}
	sm.PlrNms = map[string]struct{}{}

	// holds all team ids
	sm.TeamIds = map[string]uint64{}

	sm.SznTmPlrIds = map[int]map[uint64]map[uint64]string{}

	// setup nested team maps
	if err := sm.MapTeamIds(db); err != nil {
		fmt.Println(err)
	}

	// setup nested season maps
	if err := sm.MapSeasons(db); err != nil {
		fmt.Println(err)
	}
	return &sm
}

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapTeamIds(db *sql.DB) error {
	fmt.Println("mapping team id strings to uint64")
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
func (sm *StMaps) MapSznTeams(db *sql.DB, szn int) error {
	fmt.Println("mapping team ids to season: ", szn)
	// get all team ids
	teams, err := db.Query(
		"select distinct team_id from stats.tbox where szn_id = $1", szn)
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
		sm.SznTmPlrIds[szn][id] = map[uint64]string{}
	}
	return nil
}

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapSeasons(db *sql.DB) error {
	fmt.Println("mapping seasons")
	// to handle season id = 0
	sm.SeasonPlrNms[0] = map[string]uint64{}
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
		sm.SeasonPlrNms[szn] = map[string]uint64{}
		sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}

		if err = sm.MapSznTeams(db, szn); err != nil {
			return err
		}
	}
	return nil
}

// query all player detail from db, map to id/name
// called in UpdateStructs func to refresh the data
func (sm *StMaps) MapPlayers(db *sql.DB) error {
	fmt.Println("mapping all players")
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
		if err := sm.MapPlayerRow(db, rows, &p); err != nil {
			return fmt.Errorf("MapPlayerRow call failed\n** %w", err)
		}
	}
	return nil
}

// scan player results from db to maps
func (sm *StMaps) MapPlayerRow(db *sql.DB, rows *sql.Rows, p *StPlayer) error {
	fmt.Printf("adding %s|%d to maps\n", p.Lowr, p.Id)
	var tms string  // comma separated string to be converted to []string
	var lowr string // run remove dia on each player's lowr
	rows.Scan(&p.Id, &p.Name, &lowr, &p.Lg, &p.MaxRSzn, &p.MinRSzn,
		&p.MaxPSzn, &p.MinPSzn, &tms)

	// remove accents from lower case name
	p.Lowr = RemoveDiacritics(lowr)

	// split tms string from db to slice of strings
	sm.CleanTeamsSlice(p, tms)

	// iterate through each season played add the player to the season players map
	if err := sm.MapPlrToSzn(p); err != nil {
		return err
	}

	// for each season, query player & season and map each occurence in [szn][tm]
	if err := sm.MapSznTmPlr(db, p); err != nil {
		fmt.Printf("error occured mapping player season/teams %s | %d\n%v", p.Lowr, p.Id, err)
	}

	// basic player id and name maps for simple Exists funcs
	sm.MapPlrIdsNms(p)

	// map player struct to id & name
	sm.PlayerIdDtl[p.Id] = p
	sm.PlayerNameDtl[p.Lowr] = p

	// map id to name & name to id
	sm.PlayerIdName[p.Id] = p.Lowr
	sm.PlayerNameId[p.Lowr] = p.Id
	return nil
}

// called from within season loop
func (sm *StMaps) MapSznTmPlr(db *sql.DB, p *StPlayer) error {
	fmt.Printf("mapping %s|%d to season team maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	q := `
select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and szn_id between $2 and $3
group by player_id, szn_id`

	tmsRows, err := db.Query(q, p.Id, p.MinRSzn, p.MaxRSzn)
	if err != nil {
		return fmt.Errorf(
			"failed to query teams by season for %s|%d seasons %d through %d\n%w",
			p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn, err)
	}
	count := 0
	for tmsRows.Next() {
		count++
		fmt.Println("mapping row", count)
		var szn int
		var tmStr string
		tmsRows.Scan(&szn, &tmStr)

		tmsItr := strings.SplitSeq(tmStr, ",")

		for t := range tmsItr {
			fmt.Println("")
			teamId := sm.TeamIds[t]
			if sm.SznTmPlrIds[szn] == nil {
				sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}
			}
			// ensure inner map for team exists
			if sm.SznTmPlrIds[szn][teamId] == nil {
				sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
			}
			sm.SznTmPlrIds[szn][teamId][p.Id] = p.Lowr
		}
	}
	return nil
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
func (sm *StMaps) MapPlrToSzn(p *StPlayer) error {
	fmt.Printf("mapping %s|%d to season maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	sm.SeasonPlrNms[0][p.Lowr] = p.Id
	sm.SeasonPlrIds[0][p.Id] = p.Lowr
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrNms[int(s)][p.Lowr] = p.Id
		sm.SeasonPlrIds[int(s)][p.Id] = p.Lowr
	}
	return nil
}
