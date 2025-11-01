package api

import (
	"fmt"
	"strconv"
	"time"
)

// hold current nba and wnba seasons based on date
type SeasonLeague struct {
	SznId  uint64
	WSznId uint64
	Szn    string
	WSzn   string
}

// used in GetCurrentSzns to make seasons/ids from current year
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
func (cs *CurrentSeasons) GetCurrentSzns(now time.Time) {
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
		msg := "error converting month to int"
		fmt.Printf("%s\n%v\n", msg, err)
	}
	cs.CurSznId = cint

	// append a 2 to front of prev year, return as uint64
	pint, err := strconv.ParseUint("2"+pyyy, 10, 64)
	if err != nil {
		msg := "error converting month to int"
		fmt.Printf("%s\n%v\n", msg, err)
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
func (cs *CurrentSeasons) LgSznsByMonth(now time.Time) SeasonLeague {
	// var cs CurrentSeasons
	cs.GetCurrentSzns(now)

	// convert current month to int
	m, err := strconv.Atoi(now.Format("1"))
	if err != nil {
		msg := "error converting month to int"
		fmt.Printf("%s\n%v\n", msg, err)

	}
	fmt.Println("month: ", m)

	// convert current day to int
	d, err := strconv.Atoi(now.Format("2"))
	if err != nil {
		msg := "error converting month to int"
		fmt.Printf("%s\n%v\n", msg, err)
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
