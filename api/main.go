package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/jdetok/go-api-jdeko.me/internal/env"
	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
	"github.com/jdetok/go-api-jdeko.me/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// load environment variabels
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to get environment variables")
	}

	// configs go here - 8080 for testing, will derive real vals from environment
	cfg := config{
		addr:      env.GetString("SRV_IP"),
		cachePath: env.GetString("CACHE_PATH"),
	}

	// initialize the app with the configs
	app := &application{
		config:   cfg,
		database: mariadb.InitDB(),
		// logFile := (env.GetString("LOG_PATH") + "/test.log"),
	}

	// create array of player structs
	if app.players, err = store.GetPlayers(app.database); err != nil {
		slog.Error("failed creating players array")
	}

	// create array of season structs
	if app.seasons, err = store.GetSeasons(app.database); err != nil {
		slog.Error("failed creating seasons array")
	}

	// create array of season structs
	if app.teams, err = store.GetTeams(app.database); err != nil {
		slog.Error("failed creating teams array")
	}

	// checks if cache needs refreshed every 30 seconds, refreshes if 60 sec since last
	go store.CheckCache(app.database, &app.lastUpdate,
		&app.players, &app.seasons, &app.teams,
		30*time.Second, 300*time.Second)

	mux := app.mount()
	if err := app.run(mux); err != nil {
		slog.Error("error running server")
	}
}
