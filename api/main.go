package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api/store"
	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func main() {
	// load environment variabels
	e := errd.InitErr()

	// err := godotenv.Load()
	err := envd.LoadDotEnv()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	hostaddr, err := envd.GetEnvStr("SRV_IP")
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// CONNECT TO POSTGRES
	db, err := pgdb.PostgresConn()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// configs go here - 8080 for testing, will derive real vals from environment
	/*
		cfg := config{
			addr: hostaddr,
		}
	*/
	// initialize the app with the configs
	app := &application{
		config:   config{addr: hostaddr},
		database: db,
	}
	// create array of player structs
	if app.players, err = store.GetPlayers(app.database); err != nil {
		e.Msg = "failed creating players array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.seasons, err = store.GetSeasons(app.database); err != nil {
		e.Msg = "failed creating seasons array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.teams, err = store.GetTeams(app.database); err != nil {
		e.Msg = "failed creating teams array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// checks if store needs refreshed every 30 seconds, refreshes if 60 sec since last
	go store.UpdateStructs(app.database, &app.lastUpdate,
		&app.players, &app.seasons, &app.teams,
		30*time.Second, 300*time.Second)

	// MOUNT & RUN HTTP SERVER
	mux := app.mount()
	if err := app.run(mux); err != nil {
		e.Msg = "error mounting api/http server"
		log.Fatal(e.BuildErr(err).Error())
	}

}
