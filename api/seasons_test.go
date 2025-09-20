package api

import (
	"fmt"
	"testing"
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
