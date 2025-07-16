package mariadbconn

import (
	"fmt"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/getenv"
)

func TestInitDB(t *testing.T) {
	// if err := getenv.LoadDotEnv(); err != nil {
	// 	t.Error(err)
	// }

	dbUser := getenv.GetEnvStr("DB_USER")
	dbHost := getenv.GetEnvStr("DB_HOST")
	database := getenv.GetEnvStr("DB")
	// dbUser := "go:dbgo"
	// dbHost := "10.0.13.47:3306"
	// database := "nba"
	connStr := dbUser + "@tcp(" + dbHost + ")/" + database

	db, err := CreateDBConn(connStr)
	if err != nil {
		t.Error(err)
	}

	if err := db.Ping(); err != nil {
		t.Error(err)
		// t.Errorf("InitDB() failed: db.Ping() returned an error: %e", err)
	} else {
		fmt.Println("InitDB() test passed - successfully pinged db")
	}
}
