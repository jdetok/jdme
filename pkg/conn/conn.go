package conn

import (
	"fmt"
	"os"
)

type DBEnv struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

func Load(hostN, portN, userN, passN, dbN string) (*DBEnv, error) {
	e := &DBEnv{}
	envVars := map[string]*string{
		hostN: &e.Host,
		portN: &e.Port,
		userN: &e.User,
		passN: &e.Pass,
		dbN:   &e.Database,
	}
	for ev, v := range envVars {
		var tmp string
		if tmp = os.Getenv(ev); tmp == "" {
			return nil, fmt.Errorf("must set %s in .env", ev)
		}
		*v = tmp
	}
	return e, nil
}
