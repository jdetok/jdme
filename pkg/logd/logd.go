package logd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/mgo"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	quietLvls []string
	loudLvls  []string
	httpLvls  []string
	lg        *log.Logger
	qlg       *log.Logger
	Mongo     *mgo.MongoLogger
}

type HTTPLog struct {
	ID         bson.ObjectID `bson:"-"`
	ReqTime    time.Time     `bson:"request_time"`
	URL        string        `bson:"url"`
	Method     string        `bson:"http_method"`
	Referer    string        `bson:"referer"`
	RemoteAddr string        `bson:"remote_addr"`
	UserAgent  string        `bson:"user_agent"`
	Header     http.Header   `bson:"header"`
}

// HIGH LEVEL FUNCS TO CALL IN SOURCE
// log HTTP to mongo db server and to log file
func (l *Logd) HTTPf(r *http.Request) {
	hl := NewHTTPLog(r)
	if err := log_http_mgo(hl, l.Mongo); err != nil {
		l.log(ERROR, "error logging HTTP to mongo: %v", err)
	}
	l.log(HTTP, "served %s : %s\n", r.URL.String(), hl.ID.String())
}

// most general logging message, prints to main log file (lg)
func (l *Logd) Infof(msg string, args ...any) { l.log(INFO, msg, args...) }

// debugger level - only prints to dbg log file (qlg)
func (l *Logd) Debugf(msg string, args ...any) { l.log(DEBUG, msg, args...) }

// warnings, prints to main log file
func (l *Logd) Warnf(msg string, args ...any) { l.log(WARNING, msg, args...) }

// non fatal errors
func (l *Logd) Errorf(msg string, args ...any) { l.log(ERROR, msg, args...) }

// for graceful shutdowns - logs to main log file
func (l *Logd) Quitf(msg string, args ...any) { l.log(QUIT, msg, args...) }

// for fatal errors - mimics log.Fatalf behavior (log then exit)
func (l *Logd) Fatalf(msg string, args ...any) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

// output to io writers associated with l -> all high level Infof etc funcs call this
func (l *Logd) log(level, msg string, args ...any) {
	prefix := fmt.Sprintf("[%s] ", level)
	l.lg.SetPrefix(prefix)

	var msgf string = msg
	if len(args) > 0 && args[0] != nil {
		r, ok := args[0].(*http.Request)
		if ok {
			if err := log_http_mgo(NewHTTPLog(r), l.Mongo); err != nil {
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

// log http request in mongo
func log_http_mgo(hl *HTTPLog, ml *mgo.MongoLogger) error {
	result, err := ml.Coll.InsertOne(context.TODO(), hl)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if ok {
		hl.ID = id
		return nil
	}
	return nil
}

// make new logd object, holds io writers for log files and mongo client
func NewLogd(lo, qo io.Writer) *Logd {
	return &Logd{
		lg:        log.New(lo, "", log.LstdFlags|log.Lshortfile),
		qlg:       log.New(qo, "", log.LstdFlags|log.Lshortfile),
		quietLvls: []string{DEBUG},
		loudLvls:  []string{INFO, WARNING, ERROR, FATAL, QUIT, HTTP, HTTPERR},
		httpLvls:  []string{HTTP, HTTPERR},
	}
}

// log http requests
func NewHTTPLog(r *http.Request) *HTTPLog {
	return &HTTPLog{
		ReqTime:    time.Now(),
		URL:        r.RequestURI,
		Method:     r.Method,
		RemoteAddr: r.RemoteAddr,
		Referer:    r.Referer(),
		UserAgent:  r.UserAgent(),
		Header:     r.Header,
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

// actual err gets logged, just msg string gets sent as http errora
func (l *Logd) HTTPErr(w http.ResponseWriter, err error, code int, msg string, args ...any) {
	l.log(HTTPERR, fmt.Sprintf("%s\n**%v", msg, err), args...)
	http.Error(w, msg, code)
}
