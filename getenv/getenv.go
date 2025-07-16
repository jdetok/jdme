package getenv

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/applog"
	"github.com/joho/godotenv"
)

func LoadDotEnv() error {
	e := applog.AppErr{Process: "getenv.LoadDotEnv()"}
	if err := godotenv.Load(); err != nil {
		e.Msg = "*** FATAL: failed to load .env variabels"
		return e.BuildError(err)
	}
	return nil
}

func GetEnvStr(key string) (string, error) {
	e := applog.AppErr{Process: "GetEnvStr()"}
	val, ok := os.LookupEnv(key)
	if !ok {
		e.Msg = fmt.Sprintf("*** FATAL: couldn't key value for variable '%s'", key)
		return "", e.BuildError(errors.New("GetEnvStr() error"))
	}
	return val, nil
}

func GetEnvInt(key string) (int, error) {
	e := applog.AppErr{Process: "GetEnvInt()"}
	val, ok := os.LookupEnv(key)
	if !ok {
		e.Msg = fmt.Sprintf("*** FATAL: couldn't key value for variable '%s'", key)
		return 0, e.BuildError(errors.New("GetEnvStr() error"))
	}

	// convert key from string to int
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		e.Msg = fmt.Sprintf("*** FATAL: couldn't key value for variable '%s'", key)
		return 0, e.BuildError(errors.New("error converting to int"))
	}
	return valAsInt, nil
}
