package memd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
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

func (ms *MapStore) Persist() {
	fp := "maps.json"
	js, err := json.MarshalIndent(ms.Maps, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := os.WriteFile(fp, js, 0644); err != nil {
		fmt.Println(err)
		return
	}
}

func (ms *MapStore) BuildFromPersist() {
	fp := "maps.json"
	b, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := json.Unmarshal(b, ms.Maps); err != nil {
		fmt.Println(err)
		return
	}
}

// rewrite global stores of players, teams, etc using hash maps rather than arrays
type StMaps struct {
	// read write mutex for safe concurrent writes to map
	mu sync.RWMutex

	// string team id as key, converted uint64 team id as val
	TeamIds    map[string]uint64 // string(1610612747):uint64(1610612747)
	TeamIdLg   map[uint64]int
	TmIdAbbr   map[uint64]string // 1610612747:lal
	TmAbbrId   map[string]uint64 // lal:1610612747
	LgTmIdAbbr map[int]map[uint64]string
	LgTmAbbrId map[int]map[string]uint64

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
	SznTmPlrIds  map[int]map[uint64]map[uint64]string
	NSznTmPlrIds map[int]map[uint64]map[uint64]string
	WSznTmPlrIds map[int]map[uint64]map[uint64]string
}

// player struct - to be stored as value in a map of player id keys
type StPlayer struct {
	Id      uint64
	Name    string
	Lowr    string // name all lower case
	Lg      int
	MaxRSzn int
	MinRSzn int
	MaxPSzn int
	MinPSzn int
	Teams   []uint64 // teams player has played for
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
	sm.TeamIdLg = map[uint64]int{}
	sm.TmIdAbbr = map[uint64]string{}
	sm.TmAbbrId = map[string]uint64{}
	sm.LgTmIdAbbr = map[int]map[uint64]string{}
	sm.LgTmAbbrId = map[int]map[string]uint64{}
	sm.LgTmIdAbbr[0] = map[uint64]string{}
	sm.LgTmIdAbbr[1] = map[uint64]string{}
	sm.LgTmAbbrId[0] = map[string]uint64{}
	sm.LgTmAbbrId[1] = map[string]uint64{}

	// [szn][teamId][player]
	sm.SznTmPlrIds = map[int]map[uint64]map[uint64]string{}
	sm.NSznTmPlrIds = map[int]map[uint64]map[uint64]string{}
	sm.WSznTmPlrIds = map[int]map[uint64]map[uint64]string{}

	return &sm
}

// ran at start of runtime to setup empty maps
func (ms *MapStore) Setup(db *sql.DB, lg *logd.Logd) error {
	ms.Set(MakeMaps(db)) // empty maps
	if err := ms.Rebuild(db, lg); err != nil {
		return err
	} // map data
	// ms.Persist()
	return nil
}

func (ms *MapStore) SetupFromBuild(db *sql.DB, lg *logd.Logd) error {
	ms.Set(MakeMaps(db)) // empty maps
	ms.BuildFromPersist()

	// map data
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

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapTeamIdUints(db *sql.DB) error {
	fmt.Println("mapping team id strings to uint64")
	// get all team ids
	teams, err := db.Query(`
	select distinct a.team_id, lower(b.team), b.lg_id
	from stats.tbox a 
	join lg.team b on b.team_id = a.team_id
	`)
	if err != nil {
		return err
	}
	// convert each team id string to uint64
	for teams.Next() {
		var idstr string
		var abbr string
		var lg int
		teams.Scan(&idstr, &abbr, &lg)
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			return err
		}
		// uint64 team id mapped to string team id
		sm.TeamIds[idstr] = id

		// map teamid to league
		sm.TeamIdLg[id] = lg

		// id mapped to abbr | abbr mapped to id
		sm.TmIdAbbr[id] = abbr
		sm.TmAbbrId[abbr] = id

		// map team to league
		sm.LgTmIdAbbr[lg][id] = abbr
		sm.LgTmAbbrId[lg][abbr] = id
	}
	return nil
}

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapSeasons(db *sql.DB) error {
	fmt.Println("mapping seasons")
	// to handle season id = 0
	sm.SeasonPlrNms[0] = map[string]uint64{}
	sm.SeasonPlrIds[0] = map[uint64]string{}
	sm.SeasonPlrNms[1] = map[string]uint64{}
	sm.SeasonPlrIds[1] = map[uint64]string{}
	// get all season ids
	szns, err := db.Query("select distinct szn_id from stats.tbox")
	if err != nil {
		return err
	}
	sm.SznTmPlrIds[0] = map[uint64]map[uint64]string{}
	sm.NSznTmPlrIds[0] = map[uint64]map[uint64]string{}
	sm.WSznTmPlrIds[0] = map[uint64]map[uint64]string{}

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
		sm.NSznTmPlrIds[szn] = map[uint64]map[uint64]string{}
		sm.WSznTmPlrIds[szn] = map[uint64]map[uint64]string{}

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
		`select distinct a.team_id, b.lg_id 
		from stats.tbox a
		inner join lg.team b on b.team_id = a.team_id
		where szn_id = $1`, szn)
	if err != nil {
		return err
	}

	// convert each team id string to uint64
	for teams.Next() {
		// scan team id to a string
		var idstr, lgstr string
		if err := teams.Scan(&idstr, &lgstr); err != nil {
			return err
		}

		// get the team id as uint64
		teamId, err := sm.GetTeamIDUintCC(idstr)
		if err != nil {
			return err
		}
		if sm.SznTmPlrIds[szn] == nil {
			// !!
		}
		sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
		sm.SznTmPlrIds[0][teamId] = map[uint64]string{}
		switch lgstr {
		case "0":
			sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}
		case "1":
			sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
		}
	}
	return nil
}

func (sm *StMaps) MapTeamToSzn(szn int, teamId uint64) {
	sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
	fmt.Println(sm.TeamIdLg[teamId])
	switch sm.TeamIdLg[teamId] {
	case 0:
		fmt.Printf("%d played in nba playoff season %d\n", szn, teamId)
		sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}
	case 1:
		fmt.Printf("%d played in wnba playoff season %d\n", szn, teamId)
		sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
	}
}

// check if team had playoff games in season
func (sm *StMaps) MapPlayoffSzn(db *sql.DB, szn int, teamId uint64) (bool, error) {
	qTmSznExists := (`select exists (select team_id from stats.tbox where szn_id = $1 and team_id = $2)`)
	var exists bool
	row := db.QueryRow(qTmSznExists, szn, teamId)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	sm.MapTeamToSzn(szn, teamId)
	// sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
	// sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
	// sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}

	return true, nil
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
