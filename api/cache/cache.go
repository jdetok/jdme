package cache

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jdetok/go-api-jdeko.me/applog"
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
	e := applog.AppErr{Process: "UpdateStructs()"}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(*lastUpdate) > threshold {
			fmt.Printf("refreshing cache at %v: %v since last update\n",
				time.Now().Format("2006-01-02 15:04:05"), threshold)

			// REFRESH THE SEASONS ARRAY
			newSeasons, err := GetSeasons(db)
			if err != nil {
				e.Msg = "failed to get seasons"
				fmt.Println(e.BuildError(err))
			}
			*seasons = newSeasons

			// REFRESH THE PLAYERS ARRAY
			newPlayers, err := GetPlayers(db)
			if err != nil {
				e.Msg = "failed to get players"
				fmt.Println(e.BuildError(err))
			}
			*players = newPlayers

			// REFRESH THE SEASONS ARRAY
			newTeams, err := GetTeams(db)
			if err != nil {
				e.Msg = "failed to get teams"
				fmt.Println(e.BuildError(err))
			}
			*teams = newTeams

			updateTime := time.Now()
			*lastUpdate = updateTime
			fmt.Printf("finished refreshing cache at %v\n", updateTime)
		}
	}
}
