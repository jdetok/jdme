/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
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

	// initialize the app with the configs
	app := &api.App{
		Config:   api.Config{Addr: hostaddr},
		Database: db,
	}
	// create array of player structs
	if app.Players, err = store.GetPlayers(app.Database); err != nil {
		e.Msg = "failed creating players array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.Seasons, err = store.GetSeasons(app.Database); err != nil {
		e.Msg = "failed creating seasons array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.Teams, err = store.GetTeams(app.Database); err != nil {
		e.Msg = "failed creating teams array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// checks if store needs refreshed every 30 seconds, refreshes if 60 sec since last
	go store.UpdateStructs(app.Database, &app.LastUpdate,
		&app.Players, &app.Seasons, &app.Teams,
		30*time.Second, 300*time.Second)

	// MOUNT & RUN HTTP SERVER
	mux := app.Mount()
	if err := app.Run(mux); err != nil {
		e.Msg = "error mounting api/http server"
		log.Fatal(e.BuildErr(err).Error())
	}

}
