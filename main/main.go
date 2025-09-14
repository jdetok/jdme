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
	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func main() {
	e := errd.InitErr()
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
		Started:  0,
	}

	// update Players, Seasons, Teams in memory structs
	go api.CheckInMemStructs(app, 30*time.Second, 300*time.Second)

	// MOUNT & RUN HTTP SERVER
	mux := app.Mount()
	if err := app.Run(mux); err != nil {
		e.Msg = "error running api/http server"
		log.Fatal(e.BuildErr(err).Error())
	}

}
