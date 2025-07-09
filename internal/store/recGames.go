package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

type RecentGame struct {
	GameId string `json:"gameId"`
	GameDate string `json:"gameDate"`
	Final string `json:"final"`
	Overtime bool `json:"overtime"`
}

type RecentGames struct {
	Games []RecentGame `json:"recentGames"`
}

func MakeRgs(rows *sql.Rows) RecentGames {
	var rgs RecentGames
	for rows.Next() {
		var rg RecentGame
		rows.Scan(&rg.GameId, &rg.GameDate, &rg.Final, &rg.Overtime)
		rgs.Games = append(rgs.Games, rg)
	}
	return rgs
}

func (rgs *RecentGames) GetRecentGames(db *sql.DB) ([]byte, error){
	rows, err := db.Query(mariadb.RecentGames.Q)
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
