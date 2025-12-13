package logd

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPLog struct {
	ReqTime    time.Time   `bson:"request_time"`
	URL        string      `bson:"url"`
	Method     string      `bson:"http_method"`
	Referer    string      `bson:"referer"`
	RemoteAddr string      `bson:"remote_addr"`
	UserAgent  string      `bson:"user_agent"`
	Header     http.Header `bson:"header"`
}

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

// default logger for http requests
func (l *Logd) LogHTTP(r *http.Request) {
	var hl = NewHTTPLog(r)

	if err := l.Mongo.log(hl); err != nil {
		fmt.Println("error logging to mongo:", err)
	}

	l.log(HTTP, `
+++ REQUEST RECEIVED - %v
- Request URL: %v
- Method: %v
- Referrer: %v
- Remote Addr: %v
- User Agent: %v`,
		hl.ReqTime.Format("2006-01-02 15:04:05"),
		hl.URL,
		hl.Method,
		hl.RemoteAddr,
		hl.Referer,
		hl.UserAgent,
	)
}

// actual err gets logged, just msg string gets sent as http errora
func (l *Logd) HTTPErr(w http.ResponseWriter, err error, code int, msg string, args ...any) {
	l.log(HTTPERR, fmt.Sprintf("%s\n**%v", msg, err), args...)
	http.Error(w, msg, code)
}
