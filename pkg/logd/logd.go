package logd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logd struct {
	quietLvls []string
	loudLvls  []string
	httpLvls  []string
	lg        *log.Logger
	qlg       *log.Logger
	Mongo     *MongoLogger
}

func NewLogd(lo, qo io.Writer) *Logd {
	return &Logd{
		lg:        log.New(lo, "", log.LstdFlags|log.Lshortfile),
		qlg:       log.New(qo, "", log.LstdFlags|log.Lshortfile),
		quietLvls: []string{DEBUG},
		loudLvls:  []string{INFO, WARNING, ERROR, FATAL, QUIT, HTTP, HTTPERR},
		httpLvls:  []string{HTTP, HTTPERR},
	}
}

// create and return the log file
func SetupLogdF(pathfile string) (*os.File, error) {
	ts := time.Now().Format("01022006_150405")
	fname := fmt.Sprintf("%s_%s.log", pathfile, ts)
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to create file at %s\n**%w", fname, err)
	}
	return f, nil
}
