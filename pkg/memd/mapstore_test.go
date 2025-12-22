package memd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/joho/godotenv"
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
	l := logd.NewLogd(os.Stdout, os.Stdout)

	if err := godotenv.Load("../../.env"); err != nil {
		t.Error(err)
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(err)
	}

	if err := m.Rebuild(context.Background(), db, l, false); err != nil {
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
	if err := godotenv.Load("../../.env"); err != nil {
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
