package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
)

// TODO: season map[uint64]map[uint64]string

// rewrite global stores of players, teams, etc using hash maps rather than arrays
type StMaps struct {
	PlayerIdDtl   map[uint64]*StPlayer // id as key, struct as val
	PlayerNameDtl map[string]*StPlayer // mame as key, struct as val
	PlayerIdName  map[uint64]string    // id as key, name as val
	PlayerNameId  map[string]uint64    // player name as key, id as val (for lookup)
	SeasonPlayers map[uint64]map[uint64]string
	Teams         map[uint64]*StTeam
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
	Seasons map[uint64]bool
}

func (sm *StMaps) MakeMaps() {
	sm.PlayerIdDtl = map[uint64]*StPlayer{}
	sm.PlayerNameDtl = map[string]*StPlayer{}
	sm.PlayerIdName = map[uint64]string{}
	sm.PlayerNameId = map[string]uint64{}
	sm.SeasonPlayers = map[uint64]map[uint64]string{}
}

// team struct
// todo: hold old versions as well
type StTeam struct {
	Id   uint64
	Name string
	Lg   string
}

// TODO: query player, map each player's id to a StPlayer
func (sm *StMaps) MapPlayers(db *sql.DB) error {
	rows, err := db.Query(pgdb.QPlayerStore)
	if err != nil {
		return err
	}
	for rows.Next() {
		var p StPlayer
		var tms string // comma separated string to be converted to []string
		rows.Scan(&p.Id, &p.Name, &p.Lowr, &p.Lg, &p.MaxRSzn, &p.MinRSzn,
			&p.MaxPSzn, &p.MinPSzn, &tms)
		teamsStrArr := strings.SplitSeq(tms, ",")
		for t := range teamsStrArr {
			teamId, err := strconv.ParseUint(t, 10, 64)
			if err != nil {
				fmt.Printf("Error converting %s to string", t)
				continue
			}
			p.Teams = append(p.Teams, teamId)
		}
		p.Seasons = map[uint64]bool{}

		// add a season to the p.Seasons map for each season played
		for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
			p.Seasons[s] = true
			if sm.SeasonPlayers[s] == nil {
				sm.SeasonPlayers[s] = map[uint64]string{}
			}
			sm.SeasonPlayers[s][p.Id] = p.Name
		}

		sm.PlayerIdDtl[p.Id] = &p
		sm.PlayerNameDtl[p.Lowr] = &p
		sm.PlayerIdName[p.Id] = p.Lowr
		sm.PlayerNameId[p.Lowr] = p.Id
	}
	return nil
}
func (sm *StMaps) PlayedInSzn(searchP string, szn uint64) bool {
	p, ok := sm.PlayerNameDtl[searchP]
	if !ok {
		fmt.Println(searchP, "not found")
		return false
	}
	_, ok = p.Seasons[szn]
	if !ok {
		fmt.Printf("found %s, but they did not play in %d\n", p.Name, szn)
		return false
	}
	fmt.Printf("%s played in %d\n", p.Name, szn)
	return true
}
