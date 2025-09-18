package api

import (
	"testing"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func TestQueryTopLgPlayers(t *testing.T) {
	e := errd.InitErr()
	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	db, err := pgdb.PostgresConn()
	if err != nil {
		t.Error(e.BuildErr(err).Error())
	}
	QueryTopLgPlayers(db)
}
