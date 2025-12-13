package logd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LogHTTP struct {
	ReqTime    time.Time `bson:"request_time"`
	URL        string    `bson:"url"`
	Method     string    `bson:"http_method"`
	Referer    string    `bson:"referer"`
	RemoteAddr string    `bson:"remote_addr"`
	UserAgent  string    `bson:"user_agent"`
}

func SetupMongoLog() (*mongo.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var uri, host, port, user, pass, db string
	envVars := map[string]*string{
		"MONGO_HOST":                 &host,
		"MONGO_PORT":                 &port,
		"MONGO_INITDB_ROOT_USERNAME": &user,
		"MONGO_INITDB_ROOT_PASSWORD": &pass,
		"MONGO_INITDB_DATABASE":      &db,
	}

	for ev, v := range envVars {
		var tmp string
		if tmp = os.Getenv(ev); tmp == "" {
			return nil, fmt.Errorf("must set %s in .env", ev)
		}
		*v = tmp
	}

	uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		user, pass, host, port, db)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mogno at %s: %v", uri, err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{
		{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	return client, nil
}

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
