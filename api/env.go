package api

import (
	"fmt"
	"os"

	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
)

type Env struct {
	PGEnv    *conn.DBEnv
	MongoEnv *conn.DBEnv
	SrvIP    string
}

// capture environment variables ingested in container from .env file
// no need to load .env file, container loads them on build
func (e *Env) Load() error {
	pe, err := conn.Load("PG_HOST", "PG_PORT", "PG_USER_API", "PG_PASS_API", "PG_DB")
	if err != nil {
		return fmt.Errorf("failed to get postgres env: %v", err)
	}
	me, err := conn.Load("MONGO_HOST", "MONGO_PORT",
		"MONGO_INITDB_ROOT_USERNAME", "MONGO_INITDB_ROOT_PASSWORD", "MONGO_INITDB_DATABASE")
	if err != nil {
		return fmt.Errorf("failed to get mongodb env: %v", err)
	}
	ip := os.Getenv("SRV_IP")
	if ip == "" {
		return fmt.Errorf("error getting SRV_IP in env: %v", ip)
	}
	e.PGEnv = pe
	e.MongoEnv = me
	e.SrvIP = ip
	return nil
}
