package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
)

type Top5 struct {
	PlayerId uint64 `json:"player_id"`
	Player   string `json:"player"`
	Team     string `json:"team"`
	Points   uint32 `json:"points"`
}

type LgTop5 struct {
	LeagueTop5 []Top5 `json:"league"`
}

// use the LgTop5 query to get top 5 players per league
func QueryTopLgPlayers(db *sql.DB) {
	e := errd.InitErr()

	var lt LgTop5

	// current seasons by league
	sl := LgSznsByMonth()

	// query appropriate season for each league
	var lgs = [2]string{"nba", "wnba"}
	for _, lg := range lgs {
		// get appropriate season
		var sId string
		switch lg {
		case "nba":
			sId = strconv.FormatUint(sl.SznId, 10)
		case "wnba":
			sId = strconv.FormatUint(sl.WSznId, 10)
		}
		// query database
		r, err := db.Query(pgdb.LgTop5.Q, sId, lg)
		if err != nil {
			e.Msg = fmt.Sprintf(
				"failed to query database for top 5 lg players: sznId: %s | lg: %s\n",
				sId, lg)
			e.NewErr()
		}
		for r.Next() {
			var t Top5
			r.Scan(&t.PlayerId, &t.Player, &t.Team, &t.Points)
			lt.LeagueTop5 = append(lt.LeagueTop5, t)
		}
	}
	fmt.Println(lt)
}
