package pgdb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Ping() error
	Close() error
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// CONNECTION TO POSTGRES SERVER: MIGRATED TO POSTGRES FROM MARIADB 08/06/2025
/*
configs must be setup in .env file at project root
*/
func PostgresConn() (DB, error) {
	pg, err := GetEnvPG()
	if err != nil {
		return nil, fmt.Errorf("failed to get env for db: %v", err)
	}
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		msg := "error connecting to postgres"
		return nil, fmt.Errorf("%s\n%w", msg, err)
	}

	// set max connections
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(55 * time.Minute)
	return db, nil
}
