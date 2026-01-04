package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"github.com/jdetok/go-api-jdeko.me/pkg/resp"
	"github.com/joho/godotenv"
)

func TestHandleSeasonId(t *testing.T) {
	var p = memd.Player{}
	p.League = "NBA"
	p.Name = "LeBron James"
	p.PlayerId = 2544
	p.PSeasonIdMax = 42024
	p.SeasonIdMax = 22024
	p.PSeasonIdMin = 42005
	p.SeasonIdMin = 22003
	var errStr string
	testSzn := 22025
	resSzn := resp.HandleSeasonId(testSzn, &p, false, &errStr)

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
	var cs memd.CurrentSeasons
	sl, err := cs.LgSznsByMonth(dt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("NBA SeasonID | Season:", sl.SznId, "|", sl.Szn)
	fmt.Println("WNBA SeasonID | Season:", sl.WSznId, "|", sl.WSzn)
}
func TestQueryTopLgPlayers(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	db, err := pgdb.NewPGConn(&conn.DBEnv{}, &pgdb.DBConfig{})
	if err != nil {
		t.Error(err)
	}
	var cs memd.CurrentSeasons
	lt, err := memd.QueryTopLgPlayers(db, &cs, "10")
	if err != nil {
		t.Error(err)
	}
	js, err := memd.MarshalTopPlayers(&lt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(js))
}

// RETURN DATABASE FOR TESTING
func StartupTest(t *testing.T) pgdb.DB {
	if err := godotenv.Load("../.env"); err != nil {
		msg := "failed loading .env file"
		t.Fatal(fmt.Errorf("%s\n%v", msg, err))
	}

	db, err := pgdb.NewPGConn(&conn.DBEnv{}, &pgdb.DBConfig{})
	if err != nil {
		msg := "failed connecting to postgres"
		t.Fatal(fmt.Errorf("%s\n%v", msg, err))
	}
	return db
}

// TEST PLAYER DASH PLAYER AND TEAM TOP SCORER QUERIES
func TestGetPlayerDash(t *testing.T) {
	db := StartupTest(t)
	var pIds = []uint64{2544, 2544}    // lebron
	var sIds = []int{22024, 22024}     // 2425 reg season
	var tIds = []uint64{0, 1610612743} // first should be plr query, second tm

	for i := range pIds {
		var rp resp.RespPlayerDash
		iq := resp.PQueryIds{
			PId: pIds[i],
			TId: tIds[i],
			SId: sIds[i],
		}
		msg := fmt.Sprintf("pId: %d | sId: %d | tId: %d", iq.PId, iq.SId, iq.TId)
		js, err := rp.GetPlayerDash(db, &iq)
		if err != nil {
			emsg := fmt.Sprintf("failed getting player dash\n%s", msg)
			t.Error(fmt.Errorf("%s\n%v", emsg, err))
		}
		fmt.Println(string(js))
	}
}

// TEST PLAYERS STORE QUERY
func TestPlayerStore(t *testing.T) {
	db := StartupTest(t)
	var ps []memd.Player
	rows, err := db.Query(pgdb.PlayersSeason)
	if err != nil {
		msg := "failed getting players"
		t.Error(fmt.Errorf("%s\n%v", msg, err))
	}
	for rows.Next() {
		var p memd.Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin,
			&p.PSeasonIdMax, &p.PSeasonIdMin)
		ps = append(ps, p)
	}
}

// TEST TEAMS STORE QUERY
func TestTeamStore(t *testing.T) {
	db := StartupTest(t)
	var ts []memd.Team
	rows, err := db.Query(pgdb.Teams)
	if err != nil {
		msg := "failed getting teams"
		t.Error(fmt.Errorf("%s\n%v", msg, err))
	}

	for rows.Next() {
		var t memd.Team
		rows.Scan(&t.League, &t.TeamId, &t.TeamAbbr, &t.CityTeam)
		ts = append(ts, t)
	}
	fmt.Println(ts)
}

// TEST SEASONS STORE QUERY
func TestSeasonStore(t *testing.T) {
	db := StartupTest(t)
	var sz []memd.Season
	rows, err := db.Query(pgdb.AllSeasons)
	if err != nil {
		msg := "failed getting seasons"
		t.Error(fmt.Errorf("%s\n%v", msg, err))
	}

	for rows.Next() {
		var s memd.Season
		rows.Scan(&s.SeasonId, &s.Season, &s.WSeason)
		sz = append(sz, s)
	}
	fmt.Println(sz)
}

func TestVerifyPlayerTeam(t *testing.T) {
	db := StartupTest(t)

	iq := resp.PQueryIds{
		PId: 2544,
		SId: 22024,
		TId: 1610612747,
	}

	ptVerif, err := resp.VerifyPlayerTeam(db, &iq)
	if err != nil {
		t.Error(err)
	}
	if ptVerif {
		fmt.Printf("Player %d | Team %d | Season %d ||| verified",
			iq.PId, iq.TId, iq.SId)
	} else {
		fmt.Printf("Player %d | Team %d | Season %d ||| NOT verified",
			iq.PId, iq.TId, iq.SId)
	}
}

func TestQueryPlayerTeam(t *testing.T) {
	db := StartupTest(t)

	iq := resp.PQueryIds{
		PId: 2544,
		SId: 22024,
		TId: 1610612747,
	}

	ptVerif, err := resp.VerifyPlayerTeam(db, &iq)
	if err != nil {
		t.Error(err)
	}
	if ptVerif {
		fmt.Printf("Player %d | Team %d | Season %d ||| verified",
			iq.PId, iq.TId, iq.SId)

	} else {
		fmt.Printf("Player %d | Team %d | Season %d ||| NOT verified",
			iq.PId, iq.TId, iq.SId)
	}
}

func TestGetPlayerTeamSeason(t *testing.T) {
	db := StartupTest(t)
	iq := resp.PQueryIds{
		PId: 2544,
		SId: 22024,
		TId: 1610612747,
	}
	pltmszn, err := resp.GetPlayerTeamSeason(db, &iq)
	if err != nil {
		t.Error("failed getting pltmszn")
	}
	fmt.Println(pltmszn)
}

func TestGetTeamRecords(t *testing.T) {
	db := StartupTest(t)
	var cs memd.CurrentSeasons
	team_recs, err := memd.UpdateTeamRecords(db, &cs)
	if err != nil {
		t.Error("failed getting team records", err)
	}

	js, err := memd.TeamRecordsJSON(&team_recs)
	if err != nil {
		t.Error("failed to marshal JSON", err)
	}

	fmt.Println(string(js))
}
