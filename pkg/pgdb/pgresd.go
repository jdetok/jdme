package pgdb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jdetok/jdme/pkg/conn"
	_ "github.com/lib/pq"
)

type PostGres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	ConnStr  string
}

func GetEnvPG() (*PostGres, error) {
	var pg PostGres

	envVars := map[string]*string{
		"PG_HOST": &pg.Host,
		"PG_PORT": &pg.Port,
		"PG_USER": &pg.User,
		"PG_PASS": &pg.Password,
		"PG_DB":   &pg.Database,
	}

	for ev, v := range envVars {
		var tmp string
		if tmp = os.Getenv(ev); tmp == "" {
			return nil, fmt.Errorf("must set %s in .env", ev)
		}
		*v = tmp
	}
	return &pg, nil
}

func NewPG(e *conn.DBEnv) *PostGres {
	return &PostGres{
		Host:     e.Host,
		Port:     e.Port,
		User:     e.User,
		Password: e.Pass,
		Database: e.Database,
	}
}

func (pg *PostGres) MakeConnStr() {
	pg.ConnStr = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pg.Host, pg.Port, pg.User, pg.Password, pg.Database)
}
func (pg *PostGres) Conn() (DB, error) {
	db, err := sql.Open("postgres", pg.ConnStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"error pining postgres after successful conn: %e", err)
	}
	return db, err
}
