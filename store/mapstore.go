package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"unicode"

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

	if err := temp.MapPlayersCC(db); err != nil {
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
	mu            sync.RWMutex
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
	fmt.Println("creating empty team maps")
	if err := sm.MapTeamIds(db); err != nil {
		fmt.Println(err)
	}

	// setup nested season maps
	fmt.Println("creating empty season maps")
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
