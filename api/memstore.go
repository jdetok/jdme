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
		app.LastUpdate = app.StartTime
		app.Lg.Infof("in-memory storage goroutine first started at %v", app.StartTime)

		if err := app.Store.Rebuild(app.DB); err != nil {
			return err
		}
		app.Lg.Infof("app.Store refreshed | %d players | %d seasons | %d teams",
			len(app.Store.Players), len(app.Store.Seasons), len(app.Store.Teams))

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
				app.Lg.Infof("rebuilding memory store: {%v} since {%v} | threshold {%v}",
					time.Since(app.LastUpdate), app.LastUpdate.Format("2006-01-02 15:04:05"), threshold)

				app.LastUpdate = time.Now()

				if err := app.RebuildMemStore(ctx); err != nil {
					return fmt.Errorf("RebuildMemStore failed: %v", err)
				}
			}
		}
	}
}

func (app *App) RebuildMemStore(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(50)
	g.Go(func() error {
		if err := app.Store.Rebuild(app.DB); err != nil {
			return fmt.Errorf("error updating struct slices: %w", err)
		}
		app.Lg.Infof("app.Store refreshed | %d players | %d seasons | %d teams",
			len(app.Store.Players), len(app.Store.Seasons), len(app.Store.Teams))
		return nil
	})

	// update maps
	g.Go(func() error {
		// persist = false for continuous rebuild
		if err := app.MStore.Rebuild(ctx, app.DB, app.Lg, false); err != nil {
			return fmt.Errorf("error updating maps: %w", err)
		}
		app.Lg.Infof("app.MStore refreshed | %d players | %d teams",
			len(app.MStore.Maps.PlayerIdName), len(app.MStore.Maps.TeamIds))
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
