package resp

import (
	"database/sql"
	"encoding/json"

	"github.com/jdetok/go-api-jdeko.me/applog"
	"github.com/jdetok/go-api-jdeko.me/mdb"
)

func MakeRgs(rows *sql.Rows) RecentGames {
	var rgs RecentGames
	for rows.Next() {
		var rg RecentGame
		var ps PlayerBasic
		rows.Scan(&rg.GameId, &rg.TeamId, &rg.PlayerId,
			&rg.Player, &rg.League, &rg.Team,
			&rg.TeamName, &rg.GameDate, &rg.Matchup,
			&rg.Final, &rg.Overtime, &rg.Points, &ps.Points)

		ps.PlayerId = rg.PlayerId
		ps.TeamId = rg.TeamId
		ps.Player = rg.Player
		ps.League = rg.League
		rgs.TopScorers = append(rgs.TopScorers, ps)
		rgs.Games = append(rgs.Games, rg)
	}
	return rgs
}

func (rgs *RecentGames) GetRecentGames(db *sql.DB) ([]byte, error) {
	e := applog.AppErr{Process: "GetRecentGames()"}
	rows, err := db.Query(mdb.RecentGamePlayers.Q)
	if err != nil {
		e.Msg = "query failed"
		return nil, e.BuildError(err)
	}
	recentGames := MakeRgs(rows)
	js, err := json.Marshal(recentGames)
	if err != nil {
		e.Msg = "json marshal failed"
		return nil, e.BuildError(err)
	}
	return js, nil
}
