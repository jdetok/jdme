package logd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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
func (ml *MongoLogger) log(hl *HTTPLog) error {
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

// actual err gets logged, just msg string gets sent as http errora
func (l *Logd) HTTPErr(w http.ResponseWriter, err error, code int, msg string, args ...any) {
	l.log(HTTPERR, fmt.Sprintf("%s\n**%v", msg, err), args...)
	http.Error(w, msg, code)
}
