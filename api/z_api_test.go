package api

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func TestHandleSeasonId(t *testing.T) {

	var p = Player{}
	p.League = "NBA"
	p.Name = "LeBron James"
	p.PlayerId = 2544
	p.PSeasonIdMax = 42024
	p.SeasonIdMax = 22024
	p.PSeasonIdMin = 42005
	p.SeasonIdMin = 22003
	var errStr string
	testSzn := uint64(22025)
	resSzn := HandleSeasonId(testSzn, &p, &errStr)

	if resSzn != testSzn {
		fmt.Printf("season was manipulated: test season: %d result season %d\n",
			testSzn, resSzn)
	} else {
		fmt.Printf("season validated: test %d | result %d\n", testSzn, resSzn)
	}

	if resSzn < p.SeasonIdMin || resSzn > p.PSeasonIdMax {
		t.Errorf(`resulting season out of range:
			%d resulting season
			%d min season
			%d max season
			`, resSzn, p.SeasonIdMax, p.SeasonIdMin)
	}
}
func TestLgSznsByMonth(t *testing.T) {
	tstDate := "2026-06-21"

	dt, err := time.Parse("2006-01-02", tstDate)
	if err != nil {
		fmt.Println("error making test date")
		return
	}
	fmt.Println("test date: ", dt)
	var cs CurrentSeasons
	sl := cs.LgSznsByMonth(dt)
	fmt.Println("NBA SeasonID | Season:", sl.SznId, "|", sl.Szn)
	fmt.Println("WNBA SeasonID | Season:", sl.WSznId, "|", sl.WSzn)
}
func TestQueryTopLgPlayers(t *testing.T) {
	e := errd.InitErr()
	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	var cs CurrentSeasons
	lt, err := QueryTopLgPlayers(db, &cs, "10")
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	js, err := MarshalTopPlayers(&lt)
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	fmt.Println(string(js))
}

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
	rows, err := db.Query(pgdb.PlayersSeason)
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
	// fmt.Println(ps)
}

// TEST TEAMS STORE QUERY
func TestTeamStore(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var ts []Team
	rows, err := db.Query(pgdb.Teams)
	if err != nil {
		e.Msg = "failed getting teams"
		t.Error(e.BuildErr(err))
	}

	for rows.Next() {
		var t Team
		rows.Scan(&t.League, &t.TeamId, &t.TeamAbbr, &t.CityTeam)
		ts = append(ts, t)
	}
	// fmt.Println(ts)
}

// TEST SEASONS STORE QUERY
func TestSeasonStore(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	var sz []Season
	rows, err := db.Query(pgdb.AllSeasons)
	if err != nil {
		e.Msg = "failed getting seasons"
		t.Error(e.BuildErr(err))
	}

	for rows.Next() {
		var s Season
		rows.Scan(&s.SeasonId, &s.Season, &s.WSeason)
		sz = append(sz, s)
	}
	// fmt.Println(sz)
}
