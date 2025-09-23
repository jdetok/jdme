package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jdetok/golib/errd"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Season struct {
	SeasonId string `json:"season_id"`
	Season   string `json:"season"`
	WSeason  string `json:"wseason"`
}

type Team struct {
	League   string `json:"league"`
	TeamId   string `json:"team_id"`
	TeamAbbr string `json:"team"`
	CityTeam string `json:"team_long"`
	LogoUrl  string `json:"-"`
}

type SeasonLeague struct {
	SznId  uint64
	WSznId uint64
	Szn    string
	WSzn   string
}

type CurrentSeasons struct {
	PrevSznId uint64
	CurSznId  uint64
	PrevSzn   string
	CurSzn    string
}

/*
returns slice of two season strings for date (generally pass time.Now())
calling in 2025 will return 2024-25 and 2025-26 and so on
*/
func (cs *CurrentSeasons) CurrentSzns(now time.Time, e *errd.Err) {
	// dt := time.Now()
	dt := now
	// current year | year + 1 || e.g. 2025: cyyy=2025, cy=26
	var cyyy string = dt.Format("2006")
	var cy string = dt.AddDate(1, 0, 0).Format("06")

	// year - 1 | current year || e.g. 2025: pyyy=2024, py=25
	var pyyy string = dt.AddDate(-1, 0, 0).Format("2006")
	var py string = dt.Format("06")

	// in 2025: "2024-25", "2025-26"
	cs.PrevSzn = fmt.Sprint(pyyy, "-", py)
	cs.CurSzn = fmt.Sprint(cyyy, "-", cy)

	// append a 2 to front of current year, return as uint64
	cint, err := strconv.ParseUint("2"+cyyy, 10, 64)
	if err != nil {
		e.Msg = "error converting month to int"
		fmt.Println(e.BuildErr(err))
	}
	cs.CurSznId = cint

	// append a 2 to front of prev year, return as uint64
	pint, err := strconv.ParseUint("2"+pyyy, 10, 64)
	if err != nil {
		e.Msg = "error converting month to int"
		fmt.Println(e.BuildErr(err))
	}
	cs.PrevSznId = pint
}

/*
return SeasonLeague struct with current wnba and nba season based on the current
month. for any given year there will be two season combinations that can exist be
created using only the year as an int. for example, in 2025, both "2024-25" and
"2025-26" can be generated from the year. since the WNBA season starts and ends
in the same calendar year and the NBA season spans two calendar years, there are
times of year in which the "current" WNBA season is different than the current
NBA season.
*/
func LgSznsByMonth(now time.Time) SeasonLeague {
	e := errd.InitErr()
	var cs CurrentSeasons
	cs.CurrentSzns(now, &e)

	// convert current month to int
	m, err := strconv.Atoi(now.Format("1"))
	if err != nil {
		e.Msg = "error converting month to int"
		fmt.Println(e.BuildErr(err))
	}
	fmt.Println("month: ", m)

	// convert current day to int
	d, err := strconv.Atoi(now.Format("2"))
	if err != nil {
		e.Msg = "error converting month to int"
		fmt.Println(e.BuildErr(err))
	}
	fmt.Println("day: ", d)

	// set prev year at first (jan - april)
	var sl = SeasonLeague{
		SznId:  cs.PrevSznId,
		Szn:    cs.PrevSzn,
		WSznId: cs.PrevSznId,
		WSzn:   cs.PrevSzn,
	}

	// may through september - WNBA gets current szn, NBA gets previous
	if m > 5 {
		sl.WSznId = cs.CurSznId
		sl.WSzn = cs.CurSzn
	}

	// october 21 through end of year - both leagues get current szn
	// this is based on the 2025-26 NBA season starting on 10/21 - update day each year
	if m >= 10 && d >= 21 {
		sl.SznId = cs.CurSznId
		sl.Szn = cs.CurSzn
		sl.WSznId = cs.CurSznId
		sl.WSzn = cs.CurSzn
	}

	return sl
}

/*
use the transform package to remove accidentals
e.g. Dončić becomes doncic
*/
func RemoveDiacritics(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	output, _, _ := transform.String(t, input)
	return output
}

/*
use league and team id to generate URL with team's logo
*/
func (t Team) MakeLogoUrl() string {
	lg := strings.ToLower(t.League)
	return ("https://cdn." + lg + ".com/logos/" +
		lg + "/" + t.TeamId + "/primary/L/logo.svg")
}
