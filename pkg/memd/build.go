package memd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jdetok/jdme/pkg/errd"
	"github.com/jdetok/jdme/pkg/logd"
	"github.com/jdetok/jdme/pkg/pgdb"
)

// INITIAL MAP SETUP: must create empty maps before attempting to insert keys
// calls MapTeams and MapSeasons to setup an empty map for
// each season and team nested map
func MakeMaps() *StMaps {
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

// assign pointer to rebuilt maps struct to existing struct
func (ms *MapStore) Set(newMaps *StMaps) {
	ms.mu.Lock()
	ms.Maps = newMaps
	ms.mu.Unlock()
}

// rebuild maps in new temp StMaps structs, replace old one
func (ms *MapStore) Rebuild(ctx context.Context, db pgdb.DB, lg *logd.Logd, persist bool) error {
	temp := MakeMaps()

	if err := temp.MapTeamIdUints(db); err != nil {
		return fmt.Errorf("error mapping teams: %v", err)
	}
	if err := temp.MapSeasons(db); err != nil {
		return fmt.Errorf("error mapping seasons: %v", err)
	}
	if err := temp.MapPlayersCC(ctx, db, lg); err != nil {
		return fmt.Errorf("error mapping players: %v", err)
	}

	ms.Set(temp)

	if persist {
		if err := ms.Persist(true); err != nil {
			return &errd.PersistError{
				Err: fmt.Errorf("failed to persist memory to %s: %v",
					ms.PersistPath, err)}
		}
	}

	return nil
}

func (ms *MapStore) BuildFromPersist() error {
	b, err := os.ReadFile(ms.PersistPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, ms.Maps); err != nil {
		return err
	}
	return nil
}

// marshal to a json file for quick reload
func (ms *MapStore) Persist(mini bool) error {
	// fp := "maps.json"
	var js []byte
	var err error
	if mini {
		js, err = json.Marshal(ms.Maps)
	} else {
		js, err = json.MarshalIndent(ms.Maps, "", " ")
	}
	if err != nil {
		return err
	}
	return os.WriteFile(ms.PersistPath, js, 0644)
}
