package getenv

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/apperr"
	"github.com/joho/godotenv"
)

func LoadDotEnv() error {
	e := apperr.AppErr{Process: "getenv.LoadDotEnv()"}
	if err := godotenv.Load("../.env"); err != nil {
		e.Msg = "*** FATAL: failed to load .env variabels"
		return e.BuildError(err)
	}
	return nil
}

func GetEnvStr(key string) (string, error) {
	e := apperr.AppErr{Process: "GetEnvStr()"}
	val, ok := os.LookupEnv(key)
	if !ok {
		e.Msg = fmt.Sprintf("*** FATAL: couldn't key value for variable '%s'", key)
		return "", e.BuildError(errors.New("GetEnvStr() error"))
	}
	return val, nil
}

// func GetEnvStr(key string) string {
// 	LoadDotEnv()
// 	val, found := os.LookupEnv(key)
// 	if !found {
// 		return ""
// 	}
// 	return val
// }

func GetEnvInt(key string) int {
	LoadDotEnv()
	val, found := os.LookupEnv(key)
	if !found {
		return 0
	}

	// convert key from string to int
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return valAsInt
}
