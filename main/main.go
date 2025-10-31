/*
MIGRATION FROM MARIADB DATABASE SERVER TO POSTGRES DATABASE SERVER 08/06/2025
ALSO SWITCHED FROM PROJECT-SPECIFIC LOGGER, ERROR HANDLING TO PACKAGES IN
github.com/jdetok/golib
*/

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jdetok/go-api-jdeko.me/api"
	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/go-api-jdeko.me/store"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func main() {
	e := errd.InitErr()

	// load environment variables from .env file
	err := envd.LoadDotEnv()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// set the server IP address
	hostaddr, err := envd.GetEnvStr("SRV_IP")
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// connect to bball postgres database
	db, err := pgdb.PostgresConn()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// global pointer to configs, database, etf
	app := &api.App{
		Config:   api.Config{Addr: hostaddr},
		Database: db,
		Started:  0,
		WG:       &sync.WaitGroup{},
	}
	app.MStore.Maps = &store.StMaps{}

	// build new map store
	app.MStore.Setup(app.Database)
	if err := app.MStore.Rebuild(app.Database); err != nil {
		log.Fatal(e.BuildErr(err).Error())
	}
	fmt.Println("finished setting up map store")

	// update Players, Seasons, Teams in memory structs
	go api.CheckInMemStructs(app, 30*time.Second, 300*time.Second)

	// MOUNT & RUN HTTP SERVER
	mux := app.Mount()
	if err := app.Run(mux); err != nil {
		e.Msg = "error running api/http server"
		log.Fatal(e.BuildErr(err).Error())
	}
}
