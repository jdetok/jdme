package logd

import (
	"fmt"
	"net/http"
	"os"
	"slices"
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

// HIGH LEVEL FUNCS TO CALL IN SOURCE
func (l *Logd) HTTPf(r *http.Request) {
	hl := NewHTTPLog(r)
	l.Mongo.log(hl)
	l.log(HTTP, "http request\n\t- %s\n\t- %s", hl.ID.String(), r.URL.String())
}
func (l *Logd) Infof(msg string, args ...any)  { l.log(INFO, msg, args...) }
func (l *Logd) Debugf(msg string, args ...any) { l.log(DEBUG, msg, args...) }
func (l *Logd) Warnf(msg string, args ...any)  { l.log(WARNING, msg, args...) }
func (l *Logd) Errorf(msg string, args ...any) { l.log(ERROR, msg, args...) }
func (l *Logd) Quitf(msg string, args ...any)  { l.log(QUIT, msg, args...) }
func (l *Logd) Fatalf(msg string, args ...any) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

func (l *Logd) log(level, msg string, args ...any) {
	prefix := fmt.Sprintf("[%s] ", level)
	l.lg.SetPrefix(prefix)

	var msgf string
	if len(args) > 0 && args[0] != nil {
		r, ok := args[0].(*http.Request)
		if ok {
			if err := l.Mongo.log(NewHTTPLog(r)); err != nil {
				l.lg.Printf("failed to output log msg %s", msgf)
			}
			msgf = fmt.Sprintf("http request at endoint: %s", r.URL.String())
		} else {
			msgf = fmt.Sprintf(msg, args...)
		}
	}

	if slices.Contains(l.quietLvls, level) {
		if err := l.qlg.Output(3, msgf); err != nil {
			l.lg.Printf("failed to output log msg %s", msgf)
		}
	}

	if slices.Contains(l.loudLvls, level) || slices.Contains(l.httpLvls, level) {
		if err := l.lg.Output(3, msgf); err != nil {
			l.lg.Printf("failed to output log msg %s", msgf)
		}
	}
}
