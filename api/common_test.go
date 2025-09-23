package api

import (
	"fmt"
	"testing"
	"time"
)

func TestLgSznsByMonth(t *testing.T) {
	tstDate := "2026-06-21"

	dt, err := time.Parse("2006-01-02", tstDate)
	if err != nil {
		fmt.Println("error making test date")
		return
	}
	fmt.Println("test date: ", dt)
	sl := LgSznsByMonth(dt)
	// sl := LgSznsByMonth(time.Now())
	fmt.Println("NBA SeasonID | Season:", sl.SznId, "|", sl.Szn)
	fmt.Println("WNBA SeasonID | Season:", sl.WSznId, "|", sl.WSzn)
}
