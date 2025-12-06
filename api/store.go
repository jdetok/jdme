package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
)

// check whether enough time has passed to rebuild the in memory storage
func (app *App) UpdateStore(quickstart bool, threshold time.Duration) error {
	// call update func on intial run
	if !app.Started {
		if err := app.UpdateStructsSafe(); err != nil {
			return err
		}
		if quickstart {
			if err := app.MStore.SetupFromPersist(); err != nil {
				return fmt.Errorf(
					"** error failed to setup maps from persist file %s\n * %v",
					app.MStore.PersistPath, err,
				)
			}
		} else {
			if err := app.MStore.Setup(app.DB, app.Lg); err != nil {
				return fmt.Errorf("** error failed to setup maps\n * %v", err)
			}
		}
		app.Started = true
		app.Lg.Infof("finished with map store setup")
	}
	var errRtn error = nil
	// update structs every interval
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		if time.Since(app.LastUpdate) < threshold {
			// return nil // not time to update
			continue
		}
		app.Lg.Infof("refreshing in-mem store")
		app.Store.CurrentSzns.GetCurrentSzns(time.Now())

		var wg = &sync.WaitGroup{}
		wg.Add(2) // struct update producer, map update producer, error consumer

		var errs []error
		errCh := make(chan error)

		go func(wg *sync.WaitGroup, app *App) {
			defer wg.Done()
			if err := app.UpdateStructsSafe(); err != nil {
				errCh <- fmt.Errorf("error updating struct slices: %v", err)
			}
		}(wg, app)

		// update maps
		go func(wg *sync.WaitGroup, app *App) {
			defer wg.Done()
			if err := app.MStore.Rebuild(app.DB, app.Lg); err != nil {
				errCh <- fmt.Errorf("error updating maps: %v", err)
			}
		}(wg, app)

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			errs = append(errs, err)
		}

		if err := app.MStore.Persist(); err != nil {
			errs = append(errs, err)
		}

		numErrs := len(errs)
		if numErrs > 0 {
			errMsg := fmt.Sprintf("** %d error occured updating store\n")
			for i, err := range errs {
				errMsg = fmt.Sprintf("%s%s", errMsg,
					fmt.Sprintf("* error %d: %v\n", i, err),
				)
			}
			errRtn = fmt.Errorf("%s", errMsg)
		}
	}
	if errRtn != nil {
		return errRtn
	}
	return nil
}

// update players, seasons, and teams in memory structs slices
func (app *App) UpdateStructsSafe() error {
	app.Lg.Infof("updating in memory structs")

	var errP error
	msgP := "updating players structs"
	app.Store.Players, errP = memd.UpdatePlayers(app.DB)
	if errP != nil {
		return fmt.Errorf("failed %s\n%v", msgP, errP)
	}
	app.Lg.Infof("players count after update = %d", len(app.Store.Players))

	// update in memory seasons slice
	var errS error
	msgS := "updating seasons structs"
	app.Store.Seasons, errS = memd.UpdateSeasons(app.DB)
	if errS != nil {
		return fmt.Errorf("failed %s\n%v", msgS, errS)
	}

	// update in memory teams slice
	var errT error
	msgE := "updating teams structs"
	app.Store.Teams, errT = memd.UpdateTeams(app.DB)
	if errT != nil {
		return fmt.Errorf("failed %s\n%v", msgE, errP)
	}
	// update team records
	var errTR error
	msgTR := "updating team records"
	app.Store.TeamRecs, errTR = memd.UpdateTeamRecords(app.DB, &app.Store.CurrentSzns)
	if errTR != nil {
		return fmt.Errorf("failed %s\n%v", msgTR, errTR)
	}
	// update league top players
	var errLP error
	msgLP := "updating league top players struct"
	app.Store.TopLgPlayers, errLP = memd.QueryTopLgPlayers(app.DB, &app.Store.CurrentSzns, "50")
	if errLP != nil {
		return fmt.Errorf("failed %s\n%v", msgLP, errLP)
	}

	// update last update time
	updateTime := time.Now()
	app.LastUpdate = updateTime
	app.Lg.Infof("finished refreshing store")
	return nil
}
