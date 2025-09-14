package api

import (
	"fmt"
	"testing"
)

func TestLgSznsByMonth(t *testing.T) {
	sl := LgSznsByMonth()
	fmt.Println("NBA SeasonID | Season:", sl.SznId, "|", sl.Szn)
	fmt.Println("WNBA SeasonID | Season:", sl.WSznId, "|", sl.WSzn)
}
