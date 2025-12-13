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

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	hlg   *log.Logger
	Mongo *mongo.Client
	// jlg         *log.Logger
	HTTPLogJSON io.Writer
	quietLvls   []string
	loudLvls    []string
	httpLvls    []string
}

type LogHTTP struct {
	ReqTime    time.Time `bson:"request_time"`
	URL        string    `bson:"url"`
	Method     string    `bson:"http_method"`
	Referer    string    `bson:"referer"`
	RemoteAddr string    `bson:"remote_addr"`
	UserAgent  string    `bson:"user_agent"`
}

func NewLogd(lo, qo, ho io.Writer) *Logd {

	// lw := LogdWInit(&LogdW{out: lo, loudOut: lo, quietOut: qo})
	return &Logd{
		lg:  log.New(lo, "", log.LstdFlags|log.Lshortfile),
		qlg: log.New(qo, "", log.LstdFlags|log.Lshortfile),
		hlg: log.New(ho, "", log.LstdFlags|log.Lshortfile),
		// jlg:       log.New(jo, "", log.LstdFlags|log.Lshortfile),
		quietLvls: []string{DEBUG},
		loudLvls:  []string{INFO, WARNING, ERROR, FATAL, QUIT},
		httpLvls:  []string{HTTP, HTTPERR},
	}
}

func SetupMongoLog() (*mongo.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		return nil, fmt.Errorf("must set MONGODB_URI in .env")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mogno at %s: %v", uri, err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
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

// create json file for logging http
func SetupJSONLogdF(pathfile string) (*os.File, error) {
	ts := time.Now().Format("01022006_150405")
	fname := fmt.Sprintf("%s_%s.json", pathfile, ts)
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
func (l *Logd) Fatalf(msg string, args ...any) { l.log(FATAL, msg, args...); os.Exit(1) }
func (l *Logd) Quitf(msg string, args ...any)  { l.log(QUIT, msg, args...) }

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

	coll := l.Mongo.Database("httplog").Collection("log")

	var lh = LogHTTP{
		ReqTime:    time.Now(),
		URL:        r.RequestURI,
		Method:     r.Method,
		RemoteAddr: r.RemoteAddr,
		Referer:    r.Referer(),
		UserAgent:  r.UserAgent(),
	}

	result, err := coll.InsertOne(context.TODO(), lh)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)
}

// actual err gets logged, just msg string gets sent as http errora
func (l *Logd) HTTPErr(w http.ResponseWriter, err error, code int, msg string, args ...any) {
	l.log(HTTPERR, fmt.Sprintf("%s\n**%v", msg, err), args...)
	http.Error(w, msg, code)
}
