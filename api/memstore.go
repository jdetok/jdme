package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/errd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"golang.org/x/sync/errgroup"
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
					app.Lg.Infof("no persist file found at %s, skipping quickstart and building maps in background...\n* %v",
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
			app.Lg.Infof("UpdateStore exiting, context canceled")
			return nil
		case <-ticker.C:
			if time.Since(app.LastUpdate) >= threshold {
				app.Lg.Infof("time since last update {%v} > threshold {%v} - rebuilding memory store",
					time.Since(app.LastUpdate), threshold)

				app.LastUpdate = time.Now()

				if err := app.RebuildMemStore(ctx); err != nil {
					return fmt.Errorf("RebuildMemStore failed: %v", err)
				}
			}
		}
	}
}

func (app *App) RebuildMemStore(ctx context.Context) error {
	// var wg = &sync.WaitGroup{}

	// var errs []error
	// errCh := make(chan error, 2)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := app.UpdateStructsSafe(); err != nil {
			return fmt.Errorf("error updating struct slices: %w", err)
		}
		return nil
	})

	// update maps
	g.Go(func() error {
		// persist = false for continuous rebuild
		if err := app.MStore.Rebuild(ctx, app.DB, app.Lg, false); err != nil {
			return fmt.Errorf("error updating maps: %w", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error in RebuildMemStore: %v", err)
	}

	if err := app.MStore.Persist(true); err != nil {
		return fmt.Errorf("error persisting maps: %v", &errd.PersistError{Err: err})
	}
	return nil
}

// update players, seasons, and teams in memory structs slices
// should be in memd
func (app *App) UpdateStructsSafe() error {
	var err error
	app.Store.Players, err = memd.UpdatePlayers(app.DB)
	if err != nil {
		return fmt.Errorf("failed updating players structs: %v", err)
	}
	app.Store.Seasons, err = memd.UpdateSeasons(app.DB)
	if err != nil {
		return fmt.Errorf("failed updating season structs: %v", err)
	}
	app.Store.Teams, err = memd.UpdateTeams(app.DB)
	if err != nil {
		return fmt.Errorf("failed updating team structs: %v", err)
	}

	app.Store.TeamRecs, err = memd.UpdateTeamRecords(app.DB, &app.Store.CurrentSzns)
	if err != nil {
		return fmt.Errorf("failed updating team records: %v", err)
	}
	app.Store.TopLgPlayers, err = memd.QueryTopLgPlayers(app.DB, &app.Store.CurrentSzns, "50")
	if err != nil {
		return fmt.Errorf("failed updating league top players struct: %v", err)
	}
	app.Lg.Infof("struct slices refreshed | %d players | %d seasons | %d teams",
		len(app.Store.Players), len(app.Store.Seasons), len(app.Store.Teams))
	return nil
}
