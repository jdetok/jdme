package memd

import (
	"strconv"
)

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
	_, ok := sm.SznTmPlrIds[szn][tmId][plrId]
	return ok
}
