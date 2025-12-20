package memd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

// INITIAL MAP SETUP: must create empty maps before attempting to insert keys
// calls MapTeams and MapSeasons to setup an empty map for
// each season and team nested map
func MakeMaps() *StMaps {
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

// marshal to a json file for quick reload
func (ms *MapStore) Persist() error {
	// fp := "maps.json"
	js, err := json.MarshalIndent(ms.Maps, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(ms.PersistPath, js, 0644); err != nil {
		return err
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

// ran at start of runtime to setup empty maps
func (ms *MapStore) Setup(db pgdb.DB, lg *logd.Logd) error {
	ms.Set(MakeMaps()) // empty maps
	if err := ms.Rebuild(db, lg); err != nil {
		return err
	} // map data
	fmt.Printf("len after rebuild: %d\n", len(ms.Maps.PlrIds))
	return nil
}

func (ms *MapStore) SetupFromPersist() error {
	ms.Set(MakeMaps()) // empty maps
	if err := ms.BuildFromPersist(); err != nil {
		return err
	}
	return nil
}

// assign pointer to rebuilt maps struct to existing struct
func (ms *MapStore) Set(newMaps *StMaps) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Maps = newMaps
}

// rebuild maps in new temp StMaps structs, replace old one
func (ms *MapStore) Rebuild(db pgdb.DB, lg *logd.Logd) error {
	fmt.Println("Rebuilding StMaps...")
	temp := MakeMaps()

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

	ms.Set(temp)
	return nil
}
