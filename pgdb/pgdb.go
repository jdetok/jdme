package pgdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/pgresd"
)

func PostgresConn() (*sql.DB, error) {
	e := errd.InitErr()
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		e.Msg = "error connecting to postgres"
		return nil, e.BuildErr(err)
	}

	// set max connections
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db, nil
}
