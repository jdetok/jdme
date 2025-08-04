package pgdb

/*
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/pg"
	"github.com/jdetok/go-api-jdeko.me/applog"
)

func CreateDBConn(connStr string) (*sql.DB, error) {
	e := applog.AppErr{Process: "InitDB(): initialize database connection"}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		e.Msg = fmt.Sprintf("sql.Open() failed with connStr = %s", connStr)
		return nil, e.BuildError(err)
	}

	if err := db.Ping(); err != nil {
		e.Msg = "db.Ping() failed with returned db connection"
		return nil, e.BuildError(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	return db, nil
}
*/
