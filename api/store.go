package api

import (
	"sync"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
)

// check whether enough time has passed to rebuild the in memory storage
func (app *App) CheckInMemStructs(interval, threshold time.Duration) {
	// call update func on intial run
	if !app.Started {
		app.UpdateStructs()

		app.Started = true
	}

	// update structs every interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(app.LastUpdate) > threshold {
			app.Lg.Infof("refreshing in-mem store")

			var wg sync.WaitGroup

			// calculate current NBA/WNBA seasons
			app.Store.CurrentSzns.GetCurrentSzns(time.Now())

			// update in memory structs
			wg.Add(1)
			go func(wg *sync.WaitGroup, app *App) {
				defer wg.Done()
				app.UpdateStructs()
			}(&wg, app)

			// update maps
			wg.Add(1)
			go func(wg *sync.WaitGroup, app *App) {
				defer wg.Done()
				if err := app.MStore.Rebuild(app.DB, app.Lg); err != nil {
					app.Lg.Errorf("failed to update player map")
				}
			}(&wg, app)
			wg.Wait()

			app.MStore.Persist()
		}
	}
}

// update players, seasons, and teams in memory structs slices
func (app *App) UpdateStructs() {
	app.Lg.Infof("updating in memory structs")

	var errP error
	msgP := "updating players structs"
	app.Store.Players, errP = memd.UpdatePlayers(app.DB)
	if errP != nil {
		app.Lg.Errorf("failed %s\n%v", msgP, errP)
	}

	// update in memory seasons slice
	var errS error
	msgS := "updating seasons structs"
	app.Store.Seasons, errS = memd.UpdateSeasons(app.DB)
	if errS != nil {
		app.Lg.Errorf("failed %s\n%v", msgS, errS)
	}

	// update in memory teams slice
	var errT error
	msgE := "updating teams structs"
	app.Store.Teams, errT = memd.UpdateTeams(app.DB)
	if errT != nil {
		app.Lg.Errorf("failed %s\n%v", msgE, errP)
	}
	// update team records
	var errTR error
	msgTR := "updating team records"
	app.Store.TeamRecs, errTR = memd.UpdateTeamRecords(app.DB, &app.Store.CurrentSzns)
	if errTR != nil {
		app.Lg.Errorf("failed %s\n%v", msgTR, errTR)
	}
	// update league top players
	var errLP error
	msgLP := "updating league top players struct"
	app.Store.TopLgPlayers, errLP = memd.QueryTopLgPlayers(app.DB, &app.Store.CurrentSzns, "50")
	if errLP != nil {
		app.Lg.Errorf("failed %s\n%v", msgLP, errLP)
	}

	// update last update time
	updateTime := time.Now()
	app.LastUpdate = updateTime
	app.Lg.Infof("finished refreshing store")
}
