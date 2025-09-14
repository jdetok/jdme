package api

import (
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
)

/*
launch goroutines to update the global players, seasons, and team stores
updates at every interval
*/
func CheckInMemStructs(app *App, interval, threshold time.Duration) {
	e := errd.InitErr()

	// call update func on intial run
	if app.Started == 0 {
		UpdateStructs(app, &e)
		app.Started = 1
	}

	// update structs every interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(app.LastUpdate) > threshold {
			fmt.Printf("refreshing store at %v\n", time.Now().Format(
				"2006-01-02 15:04:05"),
			)

			UpdateStructs(app, &e)
		}
	}
}

// update players, seasons, and teams in memory structs slices
func UpdateStructs(app *App, e *errd.Err) {
	var err error

	// REFRESH THE SEASONS ARRAY
	app.Seasons, err = GetSeasons(app.Database)
	if err != nil {
		e.Msg = "failed to get seasons"
		fmt.Println(e.BuildErr(err))
	}

	// REFRESH THE PLAYERS ARRAY
	app.Players, err = GetPlayers(app.Database)
	if err != nil {
		e.Msg = "failed to get players"
		fmt.Println(e.BuildErr(err))
	}

	// REFRESH THE TEAMS ARRAY
	app.Teams, err = GetTeams(app.Database)
	if err != nil {
		e.Msg = "failed to get teams"
		fmt.Println(e.BuildErr(err))
	}

	updateTime := time.Now()
	app.LastUpdate = updateTime
	fmt.Printf("finished refreshing store at %v\n", app.LastUpdate)
}
