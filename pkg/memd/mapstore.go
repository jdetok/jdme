package memd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
)

// MapStore object assigned in the global api struct
type MapStore struct {
	Maps *StMaps

	// read write mutex for safely rewriting StMaps struct
	mu sync.RWMutex
}

// rewrite global stores of players, teams, etc using hash maps rather than arrays
type StMaps struct {
	// read write mutex for safe concurrent writes to map
	mu sync.RWMutex

	// string team id as key, converted uint64 team id as val
	TeamIds map[string]uint64

	// structs that only hold whether a player exists (id and name keys)
	PlrIds map[uint64]struct{}
	PlrNms map[string]struct{}

	// player id or name (cleaned) as key, StPlayer struct as value
	PlayerIdDtl   map[uint64]*StPlayer // id as key, struct as val
	PlayerNameDtl map[string]*StPlayer // mame as key, struct as val

	// id as key, name as val and vice versa
	PlayerIdName map[uint64]string // id as key, name as val
	PlayerNameId map[string]uint64 // player name as key, id as val (for lookup)

	// map players (id:name and name:id) to a szn id
	// used to determine whether player exists in season
	SeasonPlrNms map[int]map[string]uint64
	SeasonPlrIds map[int]map[uint64]string

	// maps a player id to a team id, which is mapped to a season
	// used to determine whether a given player played for a given team in a given season
	SznTmPlrIds map[int]map[uint64]map[uint64]string
}

// player struct - to be stored as value in a map of player id keys
type StPlayer struct {
	Id      uint64
	Name    string
	Lowr    string // name all lower case
	Lg      string
	MaxRSzn int
	MinRSzn int
	MaxPSzn int
	MinPSzn int
	Teams   []uint64 // teams player has played for
}

// ran at start of runtime to setup empty maps
func (ms *MapStore) Setup(db *sql.DB, lg *logd.Logd) error {
	ms.Set(MakeMaps(db)) // empty maps
	if err := ms.Rebuild(db, lg); err != nil {
		return err
	} // map data
	return nil
}

// assign pointer to rebuilt maps struct to existing struct
func (ms *MapStore) Set(newMaps *StMaps) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Maps = newMaps
}

// rebuild maps in new temp StMaps structs, replace old one
func (ms *MapStore) Rebuild(db *sql.DB, lg *logd.Logd) error {
	fmt.Println("Rebuilding StMaps...")
	temp := MakeMaps(db)

	// setup nested team maps
	fmt.Println("creating empty team maps")
	if err := temp.MapTeamIdUints(db); err != nil {
		fmt.Println(err)
	}

	// setup nested season maps
	fmt.Println("creating empty season maps")
	if err := temp.MapSeasons(db); err != nil {
		fmt.Println(err)
	}

	if err := temp.MapPlayersCC(db, lg); err != nil {
		fmt.Println("Error in MapPlayers:", err)
		return err
	}

	fmt.Printf("Rebuild complete: %d players mapped\n", len(temp.PlayerIdDtl))
	ms.Set(temp)
	return nil
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

	// map of seasons with nested map of player ids/names (cleaned)
	sm.SeasonPlrNms = map[int]map[string]uint64{}
	sm.SeasonPlrIds = map[int]map[uint64]string{}

	// map of player ids & player names to verify ONLY if player exists in the db
	sm.PlrIds = map[uint64]struct{}{}
	sm.PlrNms = map[string]struct{}{}

	// holds all team ids
	sm.TeamIds = map[string]uint64{}

	// [szn][teamId][player]
	sm.SznTmPlrIds = map[int]map[uint64]map[uint64]string{}

	return &sm
}

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapTeamIdUints(db *sql.DB) error {
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

// accept season as argument, query db for all teams with games played that
// season, create an empty map inside each season team map
// MapTeamIdUints MUST be called first
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
		// scan team id to a string
		var idstr string
		teams.Scan(&idstr)

		// get the team id as uint64
		teamId, err := sm.GetTeamIDUintCC(idstr)
		if err != nil {
			return err
		}

		// create an empty map (ready for player maps) inside [szn][teamId]
		sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
	}
	return nil
}

// accept comma separated string of team ids, split and append to teams slice
func (sm *StMaps) SplitTeams(p *StPlayer, tms string) ([]uint64, error) {
	var tmIds []uint64
	tmsItr := strings.SplitSeq(tms, ",")
	for t := range tmsItr {
		// access uint64 version of team id created early in sm.TeamIDs
		teamId, err := sm.GetTeamIDUintCC(t)
		if err != nil {
			return nil, err
		}
		tmIds = append(tmIds, teamId)
	}
	return tmIds, nil
}
