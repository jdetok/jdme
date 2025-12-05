package memd

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/golib/envd"
)

func TestLen(t *testing.T) {
	m := MapStore{}
	fp := "../../maps.json"
	b, err := os.ReadFile(fp)
	if err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal(b, m.Maps); err != nil {
		t.Error(err)
	}
	fmt.Println(len(m.Maps.TeamIds))
}
func TestMapPlayersCC(t *testing.T) {

	m := &MapStore{}
	l := logd.NewLogd(os.Stdout, os.Stdout, os.Stdout)

	err := envd.LoadDotEnvFile("../../.env")
	if err != nil {
		t.Error(err)
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(err)
	}

	if err := m.Rebuild(db, l); err != nil {
		t.Error(err)
	}

	var luka uint64 = 1629029
	var dal uint64 = 1610612742
	var lal uint64 = 1610612747

	var po_fail int = 42022
	var po_pass int = 42024

	p1 := m.Maps.SznTmPlrIds[po_fail][dal][luka]
	if p1 == "" {
		fmt.Printf("playoff season %d fail | %d | %d\n", po_fail, dal, luka)
	} else {
		fmt.Printf("playoff season %d pass | %d | %d | %s\n", po_fail, dal, luka, p1)
	}

	p2 := m.Maps.SznTmPlrIds[po_pass][lal][luka]
	if p2 == "" {
		fmt.Printf("playoff season %d fail | %d | %d\n", po_pass, lal, luka)
	} else {
		fmt.Printf("playoff season %d pass | %d | %d | %s\n", po_pass, lal, luka, p2)
	}

	fmt.Println(len(m.Maps.TeamIds))
}

func TestMapSznTeams(t *testing.T) {
	err := envd.LoadDotEnvFile("../../.env")
	if err != nil {
		t.Error(err)
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(err)
	}

	sm := MakeMaps()
	if err := sm.MapSeasons(db); err != nil {
		t.Error(err)
	}

	// setup nested team maps
	fmt.Println("creating empty team maps")
	if err := sm.MapTeamIdUints(db); err != nil {
		fmt.Println(err)
	}

	szn := 42024
	fmt.Println(sm.NSznTmPlrIds[szn])

	fmt.Println(sm.TeamIdLg)
}

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/jdetok/go-api-jdeko.me/api"
// 	"github.com/jdetok/go-api-jdeko.me/pgdb"
// 	"github.com/jdetok/golib/envd"
// )

// func TestMapPlayers(t *testing.T) {
// 	var sm StMaps
// 	var cs api.CurrentSeasons
// 	cs.GetCurrentSzns(time.Now())

// 	err := envd.LoadDotEnvFile("../.env")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	db, err := pgdb.PostgresConn()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	sm.MakeMaps()

// 	if err := sm.MapTeamIds(db); err != nil {
// 		t.Error(err)
// 	}

// 	if err := sm.MapPlayers(db); err != nil {
// 		t.Error(err)
// 	}

// 	testSearch := []string{"lebron james", "stephen curry", "anthony edwards"}
// 	for _, t := range testSearch {
// 		// player id from name test
// 		fmt.Printf("player search: %s | value returned: %d\n", t,
// 			sm.PlayerNameId[t])

// 		// player struct from name test
// 		fmt.Printf("player search: %s | value returned: %v\n", t,
// 			sm.PlayerNameDtl[t])
// 		plr := sm.PlayerNameDtl[t]
// 		testSzns := []uint64{22025, 22017, 22004}
// 		for _, s := range testSzns {
// 			sm.PlayedInSzn(plr.Id, s)
// 		}

// 	}
// 	fmt.Println(sm.PlayerNameDtl["eddie house"])
// 	fmt.Println(sm.SeasonPlayers[22003][2544])
// }
