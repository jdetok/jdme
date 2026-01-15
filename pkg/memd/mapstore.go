package memd

import (
	"strconv"
	"strings"
	"sync"

	"github.com/jdetok/jdme/pkg/pgdb"
)

// MapStore object assigned in the global api struct
type MapStore struct {
	Maps *StMaps

	// read write mutex for safely rewriting StMaps struct
	mu          sync.RWMutex
	PersistPath string // file path for json persist file
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

// get all team ids from db, convert each to a uint64, map to string version
func (sm *StMaps) MapTeamIdUints(db pgdb.DB) error {
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
func (sm *StMaps) MapSeasons(db pgdb.DB) error {
	sm.InitAggSznMaps([]int{0, 29999, 49999})
	szns, err := db.Query("select distinct szn_id from stats.tbox")
	if err != nil {
		return err
	}

	for szns.Next() {
		var sznStr string
		szns.Scan(&sznStr)
		szn, err := strconv.Atoi(sznStr)
		if err != nil {
			return err
		}

		sm.InitEmptySznMaps(szn)

		// adds 0 as a team id in each season map
		sm.InitTm0SznMap(szn)

		if err = sm.MapSznTeams(db, szn); err != nil {
			return err
		}
	}
	return nil
}

// create empty map for each season
func (sm *StMaps) InitEmptySznMaps(szn int) {
	sm.SeasonPlrIds[szn] = map[uint64]string{}
	sm.SeasonPlrNms[szn] = map[string]uint64{}
	sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}
	sm.NSznTmPlrIds[szn] = map[uint64]map[uint64]string{}
	sm.WSznTmPlrIds[szn] = map[uint64]map[uint64]string{}
}

func (sm *StMaps) InitAggSznMaps(aggSzns []int) {
	for _, s := range aggSzns {
		sm.SeasonPlrNms[s] = map[string]uint64{}
		sm.SeasonPlrIds[s] = map[uint64]string{}
		sm.SznTmPlrIds[s] = map[uint64]map[uint64]string{}
		sm.NSznTmPlrIds[s] = map[uint64]map[uint64]string{}
		sm.WSznTmPlrIds[s] = map[uint64]map[uint64]string{}
		sm.InitTm0SznMap(s)
		sm.InitTmSznMaps(s)
	}
}

func (sm *StMaps) InitTmSznMaps(szn int) {
	for t, lg := range sm.TeamIdLg {
		sm.SznTmPlrIds[szn][t] = map[uint64]string{}
		switch lg {
		case 0:
			sm.NSznTmPlrIds[szn][t] = map[uint64]string{}
		case 1:
			sm.WSznTmPlrIds[szn][t] = map[uint64]string{}
		}
	}
}

func (sm *StMaps) InitTm0SznMap(szn int) {
	sm.SznTmPlrIds[szn][0] = map[uint64]string{}
	sm.NSznTmPlrIds[szn][0] = map[uint64]string{}
	sm.WSznTmPlrIds[szn][0] = map[uint64]string{}
}

// accept season as argument, query db for all teams with games played that
// season, create an empty map inside each season team map
// MapTeamIdUints MUST be called first
func (sm *StMaps) MapSznTeams(db pgdb.DB, szn int) error {
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

		if sm.SznTmPlrIds[szn][teamId] == nil {
			sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
			switch lgstr {
			case "0":
				sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}
			case "1":
				sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
			}
		}

	}
	return nil
}

func (sm *StMaps) MapTeamToSzn(szn int, teamId uint64) {
	sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
	switch sm.TeamIdLg[teamId] {
	case 0:
		sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}
	case 1:
		sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
	}
}

// check if team had playoff games in season
func (sm *StMaps) MapPlayoffSzn(db pgdb.DB, szn int, teamId uint64) (bool, error) {
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
