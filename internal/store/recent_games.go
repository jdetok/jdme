package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

type RecentGames struct {
	TopScorers []PlayerBasic `json:"top_scorers"`
	Games      []RecentGame  `json:"recent_games"`
}

type PlayerBasic struct {
	PlayerId uint64 `json:"player_id"`
	TeamId   uint64 `json:"team_id"`
	Player   string `json:"player"`
	League   string `json:"league"`
	Points   uint16 `json:"points"`
}

type RecentGame struct {
	GameId   uint64 `json:"game_id"`
	TeamId   uint64 `json:"team_id"`
	PlayerId uint64 `json:"player_id"`
	Player   string `json:"player"`
	League   string `json:"league"`
	Team     string `json:"team"`
	TeamName string `json:"team_name"`
	GameDate string `json:"game_date"`
	Matchup  string `json:"matchup"`
	Final    string `json:"final"`
	Overtime bool   `json:"overtime"`
	Points   uint16 `json:"points"`
}

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
	rows, err := db.Query(mariadb.RecentGamePlayers.Q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	recentGames := MakeRgs(rows)
	js, err := json.Marshal(recentGames)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return js, nil
}
