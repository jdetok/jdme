package memd

import (
	"fmt"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/golib/envd"
)

func TestMapPlayersCC(t *testing.T) {
	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		t.Error(err)
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(err)
	}
	sm := MakeMaps(db)

	if sm.MapPlayersCC(db, &logd.Logd{}); err != nil {
		t.Error(err)
	}

	var luka uint64 = 1629029
	var dal uint64 = 1610612742
	var lal uint64 = 1610612747
	var both int = 22024
	var onlyd int = 22025
	var onlyl int = 22023

	p1 := sm.SznTmPlrIds[both][dal][luka]
	if p1 == "" {
		t.Errorf("should exist for %d | %d", both, dal)
	}
	p2 := sm.SznTmPlrIds[both][lal][luka]
	if p2 == "" {
		t.Errorf("should exist for %d | %d", both, lal)
	}
	p3 := sm.SznTmPlrIds[onlyd][dal][luka]
	if p3 == "" {
		t.Errorf("should exist for %d | %d", onlyd, dal)
	}
	p4 := sm.SznTmPlrIds[onlyd][lal][luka]
	if p4 != "" {
		t.Errorf("should NOT exist for %d | %d", onlyd, lal)

	}
	p5 := sm.SznTmPlrIds[onlyl][dal][luka]
	if p5 != "" {
		t.Errorf("should NOT exist for %d | %d", onlyl, dal)
	}
	p6 := sm.SznTmPlrIds[onlyl][lal][luka]
	if p6 == "" {
		t.Errorf("should exist for %d | %d", onlyl, lal)
	}

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

	sm := MakeMaps(db)
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
