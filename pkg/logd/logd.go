package logd

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	INFO    string = "INFO"
	DEBUG   string = "DEBUG"
	WARNING string = "* WARNING"
	ERROR   string = "** ERROR"
	FATAL   string = "*** FATAL ERROR"
	HTTP    string = "HTTP"
	HTTPERR string = "HTTPERR"
	QUIT    string = "QUIT"
)

type Logd struct {
	// lw  *LogdW
	lg    *log.Logger
	qlg   *log.Logger
	Mongo *mongo.Client
	// jlg         *log.Logger
	HTTPLogJSON io.Writer
	quietLvls   []string
	loudLvls    []string
	httpLvls    []string
}

// HIGH LEVEL FUNCS TO CALL IN SOURCE
func (l *Logd) Infof(msg string, args ...any)  { l.log(INFO, msg, args...) }
func (l *Logd) Debugf(msg string, args ...any) { l.log(DEBUG, msg, args...) }
func (l *Logd) Warnf(msg string, args ...any)  { l.log(WARNING, msg, args...) }
func (l *Logd) Errorf(msg string, args ...any) { l.log(ERROR, msg, args...) }
func (l *Logd) Quitf(msg string, args ...any)  { l.log(QUIT, msg, args...) }
func (l *Logd) Fatalf(msg string, args ...any) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

func NewLogd(lo, qo io.Writer) *Logd {
	return &Logd{
		lg:        log.New(lo, "", log.LstdFlags|log.Lshortfile),
		qlg:       log.New(qo, "", log.LstdFlags|log.Lshortfile),
		quietLvls: []string{DEBUG},
		loudLvls:  []string{INFO, WARNING, ERROR, FATAL, QUIT},
		httpLvls:  []string{HTTP, HTTPERR},
	}
}

func (l *Logd) log(level, msg string, args ...any) {
	prefix := fmt.Sprintf("[%s] ", level)
	l.lg.SetPrefix(prefix)
	msgf := fmt.Sprintf(msg, args...)

	if slices.Contains(l.quietLvls, level) {
		if err := l.qlg.Output(3, msgf); err != nil {
			l.lg.Printf("failed to output log msg %s", msgf)
		}
	}

	if slices.Contains(l.loudLvls, level) {
		if err := l.lg.Output(3, msgf); err != nil {
			l.lg.Printf("failed to output log msg %s", msgf)
		}
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
