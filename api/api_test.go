package api

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

// RETURN DATABASE FOR TESTING
func StartupTest(t *testing.T) *sql.DB {
	e := errd.InitErr()
	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		e.Msg = "failed loading .env file"
		t.Fatal(e.BuildErr(err))
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		e.Msg = "failed connecting to postgres"
		t.Fatal(e.BuildErr(err))
	}
	return db
}

// TEST PLAYER DASH PLAYER AND TEAM TOP SCORER QUERIES
func TestGetPlayerDash(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var pIds = []uint64{2544, 2544}    // lebron
	var sIds = []uint64{22024, 22024}  // 2425 reg season
	var tIds = []uint64{0, 1610612743} // first should be plr query, second tm

	for i := range pIds {
		var rp Resp
		plr := pIds[i]
		szn := sIds[i]
		tm := tIds[i]
		msg := fmt.Sprintf("pId: %d | sId: %d | tId: %d", plr, szn, tm)
		js, err := rp.GetPlayerDash(db, plr, szn, tm)
		if err != nil {
			e.Msg = fmt.Sprintf("failed getting player dash\n%s", msg)
			t.Error(e.BuildErr(err))
		}
		fmt.Println(string(js))
	}
}

// TEST PLAYERS STORE QUERY
func TestPlayerStore(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var ps []Player
	rows, err := db.Query(pgdb.PlayersSeason.Q)
	if err != nil {
		e.Msg = "failed getting players"
		t.Error(e.BuildErr(err))
	}
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin,
			&p.PSeasonIdMax, &p.PSeasonIdMin)
		ps = append(ps, p)
	}
	fmt.Println(ps)
}

// TEST TEAMS STORE QUERY
func TestTeamStore(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var ts []Team
	rows, err := db.Query(pgdb.Teams.Q)
	if err != nil {
		e.Msg = "failed getting teams"
		t.Error(e.BuildErr(err))
	}

	for rows.Next() {
		var t Team
		rows.Scan(&t.League, &t.TeamId, &t.TeamAbbr, &t.CityTeam)
		ts = append(ts, t)
	}
	fmt.Println(ts)
}

// TEST SEASONS STORE QUERY
func TestSeasonStore(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var sz []Season
	rows, err := db.Query(pgdb.AllSeasons.Q)
	if err != nil {
		e.Msg = "failed getting seasons"
		t.Error(e.BuildErr(err))
	}

	for rows.Next() {
		var s Season
		rows.Scan(&s.SeasonId, &s.Season, &s.WSeason)
		sz = append(sz, s)
	}
	fmt.Println(sz)
}
