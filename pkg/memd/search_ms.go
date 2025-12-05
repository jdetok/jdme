package memd

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/pkg/errd"
)

func (sm *StMaps) GetPlrLg(plrId uint64) (int, error) {
	p, ok := sm.PlayerIdDtl[plrId]
	if !ok {
		return 999, fmt.Errorf("couldn't find player %d", plrId)
	}
	return p.Lg, nil
}

type SznTmPlr struct {
	PId uint64
	TId uint64
	SId int
}

func (sm *StMaps) ValiSznTmPlr(plrId, tmId uint64, sId int) (*SznTmPlr, error) {
	// return if valid

	if _, ok := sm.SznTmPlrIds[sId][tmId][plrId]; ok {
		return &SznTmPlr{PId: plrId, TId: tmId, SId: sId}, nil
	}
	fmt.Println("past first one")
	// not ok, first check player exists
	if _, ok := sm.PlrIds[plrId]; !ok {
		return nil, &errd.ValidationError{Val: plrId}
		// return nil, fmt.Errorf("player %d doesn't exist", plrId)
	}

	fmt.Println("past second one")

	// check player exists in team
	if _, ok := sm.SeasonPlrIds[sId][plrId]; ok {
		return &SznTmPlr{PId: plrId, TId: 0, SId: sId}, nil
	}

	// player exists, but not in that season or team. get most recent
	if p, ok := sm.PlayerIdDtl[plrId]; ok {
		maxSzn := p.MaxRSzn
		fmt.Println("past fourth one")
		return &SznTmPlr{PId: plrId, TId: 0, SId: maxSzn}, nil
	}
	return nil, fmt.Errorf("couldn't validate %d | %d | %d", plrId, tmId, sId)
}

// THIS CAUSES ISSUE WITH NBA/WNBA ABBR OVERLAP. NEED TO TRACK LEAGUE TOO.
func (sm *StMaps) GetLgTmIdFromAbbr(abbr string, lg int) (uint64, error) {
	var tmId uint64
	var ok bool
	tmId, ok = sm.LgTmAbbrId[lg][abbr]
	if !ok {
		return 0, fmt.Errorf("couldn't get teamId from team abbr %s", abbr)
	}
	return tmId, nil
}

// return players max reg season
func (sm *StMaps) GetSznFromPlrId(plrId uint64) (int, error) {
	if p, ok := sm.PlayerIdDtl[plrId]; ok {
		fmt.Println(p)
		return p.MaxRSzn, nil
	}
	return 0, fmt.Errorf("no season found for %d", plrId)
}

// ["random"] returns 0
func (sm *StMaps) GetPlrIdFromName(name string) (uint64, error) {
	if plrId, ok := sm.PlayerNameId[name]; ok {
		return plrId, nil
	}
	return 0, fmt.Errorf("couldn't get player id from %s", name)
}

// search SeasonPlrNms or SeasonPlrIds
// move string logic to PlayerFromQ
func (sm *StMaps) PlrExistsInSzn(searchP string, searchS int) bool {
	var ok bool
	plrIdInt, err := strconv.ParseUint(searchP, 10, 64)

	// search by player id if searchP successfully converts to int
	if err == nil {
		// only verify player name exists in map if 0 season
		if searchS == 0 {
			return sm.PlrIdExists(plrIdInt)
		} // verify player id exists in passed season's map
		_, ok = sm.SeasonPlrIds[searchS][plrIdInt]
	} else {
		// search by name (lowercase, accidentals removed) if failed to convert
		// to int if season is 0
		if searchS == 0 { // only verify player id if 0 season
			return sm.PlrNmExists(searchP)
		}
		// verify player name exists in passed season's map
		_, ok = sm.SeasonPlrNms[searchS][searchP]
	}
	return ok
}

func (sm *StMaps) PlrIdSznExists(searchP uint64, searchS int) bool {
	// verify player name exists in passed season's map
	_, ok := sm.SeasonPlrIds[searchS][searchP]
	return ok
}

func (sm *StMaps) PlrSznExists(searchP string, searchS int) bool {
	// verify player name exists in passed season's map
	_, ok := sm.SeasonPlrNms[searchS][searchP]
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

// check if passed player name exists in PlrNms map
func (sm *StMaps) PlrSznTmExists(plrId, tmId uint64, szn int) bool {
	fmt.Println("in season:", len(sm.SznTmPlrIds[szn]))
	fmt.Println("in team:", len(sm.SznTmPlrIds[szn][tmId]))
	for _, p := range sm.SznTmPlrIds[szn][tmId] {
		fmt.Println(p)
	}
	_, ok := sm.SznTmPlrIds[szn][tmId][plrId]
	return ok
}

// get subset of players in tId team, sId season, lg league
// get random number <= number of those players
// return player in index of that random number
// if lg is 0 only nba, 1 only wnba, 10 both
// if tId is 0, all players in season used
func (sm *StMaps) RandomPlrIdV2(tId uint64, sId, lg int) uint64 {
	// get list of pId from [szn].values()
	var m map[uint64]string
	switch lg {
	case 0:
		m = sm.NSznTmPlrIds[sId][tId]
	case 1:
		m = sm.WSznTmPlrIds[sId][tId]
	case 10:
		m = sm.SznTmPlrIds[sId][tId]
	default:
		m = sm.SznTmPlrIds[sId][tId]
	}

	plrs := make([]uint64, 0, len(m))

	for id := range m {
		plrs = append(plrs, id)
	}

	if len(plrs) == 0 {
		return 0
	}
	randNum := rand.IntN(len(plrs))
	fmt.Println(plrs[randNum])
	return plrs[randNum]
}
