package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
)

type DBConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

func NewDBConf(maxOpen, maxIdle int, maxLife time.Duration) *DBConfig {
	return &DBConfig{
		MaxOpenConns: maxOpen,
		MaxIdleConns: maxIdle,
		ConnMaxLife:  maxLife,
	}
}

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Ping() error
	Close() error
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

func NewPGConn(e *conn.DBEnv, conf *DBConfig) (DB, error) {
	pg := NewPG(e)
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres\n%v", err)
	}
	if conf == nil {
		return db, nil
	}

	// set max connections
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLife)
	return db, nil
}

// CONNECTION TO POSTGRES SERVER: MIGRATED TO POSTGRES FROM MARIADB 08/06/2025
// configs must be setup in .env file at project root
func PostgresConn(conf *DBConfig) (DB, error) {
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

	if conf == nil {
		return db, nil
	}

	// set max connections
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLife)
	return db, nil
}
