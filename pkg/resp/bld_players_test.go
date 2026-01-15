package resp

import (
	"fmt"
	"testing"

	"github.com/jdetok/jdme/pkg/memd"
)

func TestHandleSeason(t *testing.T) {
	lbj := &memd.Player{
		PlayerId:     2544,
		Name:         "Lebron James",
		League:       "nba",
		SeasonIdMax:  22025,
		SeasonIdMin:  22003,
		PSeasonIdMax: 42024,
		PSeasonIdMin: 42005,
	}

	steph := &memd.Player{
		PlayerId:     201939,
		Name:         "Stephen Curry",
		League:       "nba",
		SeasonIdMax:  22025,
		SeasonIdMin:  22009,
		PSeasonIdMax: 42024,
		PSeasonIdMin: 42012,
	}
	cbrink := &memd.Player{
		PlayerId:     1642287,
		Name:         "Cameron Brink",
		League:       "wnba",
		SeasonIdMax:  22025,
		SeasonIdMin:  22024,
		PSeasonIdMax: 0,
		PSeasonIdMin: 0,
	}

	tests := []struct {
		desc    string
		want    int
		sId     int
		p       *memd.Player
		team    bool
		wantErr bool
	}{
		{"LeBron regular season", 22025, 22025, lbj, false, false},
		{"Curry before first playoff", 42012, 42011, steph, false, false},
		{"Cameron Brink playoffs", 22025, 42024, cbrink, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			es := ""
			retVal := HandleSeasonId(tt.sId, tt.p, tt.team, &es)
			if es != "" {
				fmt.Println("error string:", es)
			}
			fmt.Printf("returned int for %s: %d\n", tt.p.Name, retVal)
			if retVal != tt.want {
				fmt.Printf("got %d from %s for season %d | wanted %d\n",
					retVal, tt.p.Name, tt.sId, tt.want)

			}
		})
	}
}
