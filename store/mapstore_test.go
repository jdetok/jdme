package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
)

func TestMapPlayers(t *testing.T) {
	var sm StMaps
	var cs api.CurrentSeasons
	cs.GetCurrentSzns(time.Now())

	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		t.Error(err)
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(err)
	}

	sm.MakeMaps()
	if err := sm.MapPlayers(db); err != nil {
		t.Error(err)
	}

	testSearch := []string{"lebron james", "stephen curry", "anthony edwards"}
	for _, t := range testSearch {
		// player id from name test
		fmt.Printf("player search: %s | value returned: %d\n", t,
			sm.PlayerNameId[t])

		// player struct from name test
		fmt.Printf("player search: %s | value returned: %v\n", t,
			sm.PlayerNameDtl[t])
		plr := sm.PlayerNameDtl[t]
		testSzns := []uint64{22025, 22017, 22004}
		for _, s := range testSzns {
			// if s <= plr.MaxRSzn && s >= plr.MinRSzn {
			// 	fmt.Printf("%s played in %d\n", plr.Name, s)
			// } else {
			// 	fmt.Printf("%s did not play in %d\n", plr.Name, s)
			// }
			sm.PlayedInSzn(plr.Lowr, s)
		}

	}
	fmt.Println(sm.SeasonPlayers[22003])

	// fmt.Printf("player search: %s | value returned: %d", testSearch,
	// 	sm.PlayerNameId[testSearch])
}
