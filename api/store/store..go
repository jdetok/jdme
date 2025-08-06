package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
)

func UpdateStructs(
	db *sql.DB,
	lastUpdate *time.Time,
	players *[]Player,
	seasons *[]Season,
	teams *[]Team,
	interval time.Duration,
	threshold time.Duration) {

	// func starts here
	e := errd.InitErr()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(*lastUpdate) > threshold {
			fmt.Printf("refreshing store at %v: %v since last update\n",
				time.Now().Format("2006-01-02 15:04:05"), threshold)

			// REFRESH THE SEASONS ARRAY
			newSeasons, err := GetSeasons(db)
			if err != nil {
				e.Msg = "failed to get seasons"
				fmt.Println(e.BuildErr(err))
			}
			*seasons = newSeasons

			// REFRESH THE PLAYERS ARRAY
			newPlayers, err := GetPlayers(db)
			if err != nil {
				e.Msg = "failed to get players"
				fmt.Println(e.BuildErr(err))
			}
			*players = newPlayers

			// REFRESH THE SEASONS ARRAY
			newTeams, err := GetTeams(db)
			if err != nil {
				e.Msg = "failed to get teams"
				fmt.Println(e.BuildErr(err))
			}
			*teams = newTeams

			updateTime := time.Now()
			*lastUpdate = updateTime
			fmt.Printf("finished refreshing store at %v\n", updateTime)
		}
	}
}
