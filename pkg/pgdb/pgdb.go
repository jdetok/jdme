package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/golib/pgresd"
)

// CONNECTION TO POSTGRES SERVER: MIGRATED TO POSTGRES FROM MARIADB 08/06/2025
/*
configs must be setup in .env file at project root
*/
func PostgresConn() (*sql.DB, error) {
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		msg := "error connecting to postgres"
		return nil, fmt.Errorf("%s\n%w", msg, err)
	}

	// set max connections
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	return db, nil
}
