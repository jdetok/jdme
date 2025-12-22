package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
)

// check whether enough time has passed to rebuild the in memory storage
func (app *App) UpdateStore(ctx context.Context, quickstart bool, tick, threshold time.Duration) error {
	// call update func on intial run
	if !app.Started {
		app.Started = true
		app.StartTime = time.Now()

		if err := app.UpdateStructsSafe(); err != nil {
			return err
		}

		// init empty maps
		app.MStore.Set(memd.MakeMaps()) // empty maps

		if quickstart {
			// build maps from persisted JSON
			if err := app.MStore.BuildFromPersist(); err == nil {
				app.Lg.Infof("quick startup from %s complete", app.MStore.PersistPath)
			} else {
				if errors.Is(err, os.ErrNotExist) { // persist file doesn't exist, skip
					app.Lg.Infof("no persist file found at %s, skipping quickstart and building maps...\n* %v",
						app.MStore.PersistPath, err)
					if err := app.MStore.Rebuild(ctx, app.DB, app.Lg, true); err != nil {
						return fmt.Errorf("failed to build map store after skipping quickstart: %v", err)
					}
				} else {
					return fmt.Errorf("** failed to setup maps from persist file %s\n * %v",
						app.MStore.PersistPath, err,
					)
				}
			}
		} else {
			app.Lg.Infof("skipping quickstart and builiding map store")
			// persist = true for first run
			if err := app.MStore.Rebuild(ctx, app.DB, app.Lg, true); err != nil {
				return fmt.Errorf("failed to build map store or persist data after skipping quickstart: %v", err)
			}
		}
	}

	app.Lg.Infof("memstore setup complete | %d players | %d teams\n+ persisted at %s",
		len(app.MStore.Maps.PlrIds), len(app.MStore.Maps.TeamIds),
		app.MStore.PersistPath)

	app.Lg.Infof("starting ticker to update store: tick set to: %v | threshold: %v", tick, threshold)

	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			app.Lg.Infof("UpdateStore exiting: %v", ctx.Err())
			return nil
		case <-ticker.C:
			if time.Since(app.LastUpdate) >= threshold {
				app.Lg.Infof("time since last update {%v} > threshold {%v} - rebuilding memory store",
					time.Since(app.LastUpdate), threshold)
				if err := app.RebuildMemStore(ctx); err != nil {
					return fmt.Errorf("RebuildMemStore failed: %v", err)
				}
			}
		}
	}
}

func (app *App) RebuildMemStore(ctx context.Context) error {
	var wg = &sync.WaitGroup{}

	var errs []error
	errCh := make(chan error, 2)

	wg.Go(func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		if err := app.UpdateStructsSafe(); err != nil {
			select {
			case errCh <- fmt.Errorf("error updating struct slices: %w", err):
			case <-ctx.Done():
			}
		}
	})

	// update maps
	wg.Go(func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		// persist = false for continuous rebuild
		if err := app.MStore.Rebuild(ctx, app.DB, app.Lg, false); err != nil {
			select {
			case errCh <- fmt.Errorf("error updating maps: %w", err):
			case <-ctx.Done():
			}
		}
	})

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		errs = append(errs, err)
	}

	if err := app.MStore.Persist(true); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		var emsg strings.Builder
		fmt.Fprintf(&emsg, "** %d error occured updating store\n", len(errs))
		for i, err := range errs {
			fmt.Fprintf(&emsg, "* error %d: %v\n", i+1, err)
		}
		return errors.New(emsg.String())
	}
	return nil
}

// update players, seasons, and teams in memory structs slices
func (app *App) UpdateStructsSafe() error {
	var errP error
	msgP := "updating players structs"
	app.Store.Players, errP = memd.UpdatePlayers(app.DB)
	if errP != nil {
		return fmt.Errorf("failed %s\n%v", msgP, errP)
	}

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
	app.Lg.Infof(`finished refreshing in-memory struct slices
+ player count: %d | + season count: %d | + team count: %d
`, len(app.Store.Players), len(app.Store.Seasons), len(app.Store.Teams))
	return nil
}
