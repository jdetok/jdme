package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

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

type RecentGames struct {
	Games []RecentGame `json:"recent_games"`
}

func MakeRgs(rows *sql.Rows) RecentGames {
	var rgs RecentGames
	for rows.Next() {
		var rg RecentGame
		rows.Scan(&rg.GameId, &rg.TeamId, &rg.PlayerId, &rg.Player, &rg.League, &rg.Team,
			&rg.TeamName, &rg.GameDate, &rg.Matchup, &rg.Final, &rg.Overtime, &rg.Points)
		rgs.Games = append(rgs.Games, rg)
	}
	return rgs
}

func (rgs *RecentGames) GetRecentGames(db *sql.DB) ([]byte, error) {
	// rows, err := db.Query(mariadb.RecentGames.Q)
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

// func RecGamesTest() {
// 	rg := RecentGames{}
// 	db := mariadb.InitDB()
// 	rows, cols, err := mariadb.Select(db, mariadb.RecentGames.Q)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(cols)
// 	for rows.Next(){
// 		rows.Scan(&rg)
// 	}
// 	fmt.Println(rg)
// }
