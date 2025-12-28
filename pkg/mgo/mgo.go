package mgo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoAuth struct {
	conn string
	host string
	port string
	user string
	pass string
	db   string
}

type MongoEnv struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

func (e *MongoEnv) Load() error {
	envVars := map[string]*string{
		"MONGO_HOST":                 &e.Host,
		"MONGO_PORT":                 &e.Port,
		"MONGO_INITDB_ROOT_USERNAME": &e.User,
		"MONGO_INITDB_ROOT_PASSWORD": &e.Pass,
		"MONGO_INITDB_DATABASE":      &e.Database,
	}
	for ev, v := range envVars {
		var tmp string
		if tmp = os.Getenv(ev); tmp == "" {
			return fmt.Errorf("must set %s in .env", ev)
		}
		*v = tmp
	}
	return nil
}

type MongoLogger struct {
	Client *mongo.Client
	DB     *mongo.Database
	Coll   *mongo.Collection
}

func NewMongoLogger(e *conn.DBEnv, db, coll string) (*MongoLogger, error) {
	auth, err := getMongoAuth(e)
	if err != nil {
		return nil, err
	}
	mongo_cl, err := setupMongoClient(auth)
	if err != nil {
		return nil, err
	}
	mongo_db := mongo_cl.Database(db)
	mongo_coll := mongo_db.Collection(coll)
	return &MongoLogger{
		Client: mongo_cl,
		DB:     mongo_db,
		Coll:   mongo_coll,
	}, nil
}

func setupMongoClient(authstr string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(authstr)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mogno at %s: %v", authstr, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{
		{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	return client, nil
}

func getMongoAuth(e *conn.DBEnv) (string, error) {

	// envVars := map[string]*string{
	// 	"MONGO_HOST":                 &m.host,
	// 	"MONGO_PORT":                 &m.port,
	// 	"MONGO_INITDB_ROOT_USERNAME": &m.user,
	// 	"MONGO_INITDB_ROOT_PASSWORD": &m.pass,
	// 	"MONGO_INITDB_DATABASE":      &m.db,
	// }

	// for ev, v := range envVars {
	// 	var tmp string
	// 	if tmp = os.Getenv(ev); tmp == "" {
	// 		return nil, fmt.Errorf("must set %s in .env", ev)
	// 	}
	// 	*v = tmp
	// }
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin&directConnection=true",
		e.User, e.Pass, e.Host, e.Port, e.Database), nil
}
