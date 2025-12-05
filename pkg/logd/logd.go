package logd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

const (
	INFO    string = "INFO"
	DEBUG   string = "DEBUG"
	WARNING string = "* WARNING"
	ERROR   string = "** ERROR"
	FATAL   string = "*** FATAL ERROR"
	HTTP    string = "HTTP"
	HTTPERR string = "HTTPERR"
)

type Logd struct {
	// lw  *LogdW
	lg        *log.Logger
	qlg       *log.Logger
	hlg       *log.Logger
	quietLvls []string
	loudLvls  []string
	httpLvls  []string
}

func NewLogd(lo, qo, ho io.Writer) *Logd {
	// lw := LogdWInit(&LogdW{out: lo, loudOut: lo, quietOut: qo})
	return &Logd{
		lg:        log.New(lo, "", log.LstdFlags|log.Lshortfile),
		qlg:       log.New(qo, "", log.LstdFlags|log.Lshortfile),
		hlg:       log.New(ho, "", log.LstdFlags|log.Lshortfile),
		quietLvls: []string{DEBUG},
		loudLvls:  []string{INFO, WARNING, ERROR, FATAL},
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

	if slices.Contains(l.httpLvls, level) {
		if err := l.hlg.Output(3, msgf); err != nil {
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

// HIGH LEVEL FUNCS TO CALL IN SOURCE
func (l *Logd) Infof(msg string, args ...any)  { l.log(INFO, msg, args...) }
func (l *Logd) Debugf(msg string, args ...any) { l.log(DEBUG, msg, args...) }
func (l *Logd) Warnf(msg string, args ...any)  { l.log(WARNING, msg, args...) }
func (l *Logd) Errorf(msg string, args ...any) { l.log(ERROR, msg, args...) }
func (l *Logd) Fatalf(msg string, args ...any) { l.log(ERROR, msg, args...) }

// default logger for http requests
func (l *Logd) LogHTTP(r *http.Request) {
	l.log(HTTP, `
+++ REQUEST RECEIVED - %v
- Request URL: %v
- Method: %v | Request URI: %v
- Referrer: %v
- Remote Addr: %v
- User Agent: %v`,
		time.Now().Format("2006-01-02 15:04:05"),
		r.URL,
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		r.Referer(),
		r.UserAgent(),
	)
}

// actual err gets logged, just msg string gets sent as http errora
func (l *Logd) HTTPErr(w http.ResponseWriter, err error, code int, msg string, args ...any) {
	l.log(HTTPERR, fmt.Sprintf("%s\n**%v", msg, err), args...)
	http.Error(w, msg, code)
}
