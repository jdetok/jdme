package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
)

type Top5 struct {
	PlayerId uint64 `json:"player_id"`
	Player   string `json:"player"`
	Season   string `json:"season"`
	Team     string `json:"team"`
	Points   uint32 `json:"points"`
}

type LgTop5 struct {
	NBATop5  []Top5 `json:"nba"`
	WNBATop5 []Top5 `json:"wnba"`
}

// use the LgTop5 query to get top 5 players per league
/*
set a slice of strings with both leagues to loop through. NBA is first in the slice -
this must be maintained for the logic to work. at the end of the for loop the sId
variable is set to the WNBA season - it's declared as the NBA season before the loop begins
*/
func QueryTopLgPlayers(db *sql.DB, cs *CurrentSeasons, numPl string) (LgTop5, error) {
	e := errd.InitErr()

	var lt LgTop5

	// current seasons by league
	sl := cs.LgSznsByMonth(time.Now())

	// query appropriate season for each league
	var lgs = [2]string{"nba", "wnba"}
	var sId string = strconv.FormatUint(sl.SznId, 10)
	for _, lg := range lgs {
		// query database
		r, err := db.Query(pgdb.LeagueTopScorers, sId, lg, numPl)
		if err != nil {
			e.Msg = fmt.Sprintf(
				"failed to query database for top 5 lg players: sznId: %s | lg: %s\n",
				sId, lg)
			return lt, e.NewErr()
		}

		// create a Top5 struct for each row, append to appropriate NBA/WNBA member
		for r.Next() {
			var t Top5
			r.Scan(&t.PlayerId, &t.Player, &t.Season, &t.Team, &t.Points)
			switch lg {
			case "nba":
				lt.NBATop5 = append(lt.NBATop5, t)
			case "wnba":
				lt.WNBATop5 = append(lt.WNBATop5, t)
			}
		}

		// after first run set wnba season
		sId = strconv.FormatUint(sl.WSznId, 10)
	}
	return lt, nil
}

// marshal LgTop5 struct into JSON []byte
func MarshalTop5(lt *LgTop5) ([]byte, error) {
	e := errd.InitErr()
	js, err := json.Marshal(lt)
	if err != nil {
		return nil, e.BuildErr(err)
	}
	return js, nil
}
