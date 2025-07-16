package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/cache"
	"github.com/jdetok/go-api-jdeko.me/apperr"
	"github.com/jdetok/go-api-jdeko.me/getenv"
	"github.com/jdetok/go-api-jdeko.me/mdb"
)

func main() {
	// load environment variabels
	e := apperr.AppErr{Process: "Main function"}

	// err := godotenv.Load()
	err := getenv.LoadDotEnv()
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	hostaddr, err := getenv.GetEnvStr("SRV_IP")
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	cacheP, err := getenv.GetEnvStr("CACHE_PATH")
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	dbConnStr, err := getenv.GetEnvStr("DB_CONN_STR")
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	db, err := mdb.CreateDBConn(dbConnStr)
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	// configs go here - 8080 for testing, will derive real vals from environment
	cfg := config{
		addr:      hostaddr,
		cachePath: cacheP,
	}

	// initialize the app with the configs
	app := &application{
		config:   cfg,
		database: db,
	}
	// create array of player structs
	if app.players, err = cache.GetPlayers(app.database); err != nil {
		e.Msg = "failed creating players array"
		fmt.Println(e.BuildError(err).Error())
	}

	// create array of season structs
	if app.seasons, err = cache.GetSeasons(app.database); err != nil {
		e.Msg = "failed creating seasons array"
		fmt.Println(e.BuildError(err).Error())
	}

	// create array of season structs
	if app.teams, err = cache.GetTeams(app.database); err != nil {
		e.Msg = "failed creating teams array"
		fmt.Println(e.BuildError(err).Error())
	}

	// checks if cache needs refreshed every 30 seconds, refreshes if 60 sec since last
	go cache.UpdateStructs(app.database, &app.lastUpdate,
		&app.players, &app.seasons, &app.teams,
		30*time.Second, 300*time.Second)

	mux := app.mount()
	if err := app.run(mux); err != nil {
		e.Msg = "error mounting api/http server"
		log.Fatal(e.BuildError(err).Error())
	}
}
