package mgo

import (
	"context"
	"fmt"
	"os"
	"time"

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

type MongoLogger struct {
	Client *mongo.Client
	DB     *mongo.Database
	Coll   *mongo.Collection
}

func NewMongoLogger(db, coll string) (*MongoLogger, error) {
	auth, err := getMongoAuth()
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

func setupMongoClient(auth *mongoAuth) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(auth.conn)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mogno at %s: %v", auth.conn, err)
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

func getMongoAuth() (*mongoAuth, error) {
	var m mongoAuth

	envVars := map[string]*string{
		"MONGO_HOST":                 &m.host,
		"MONGO_PORT":                 &m.port,
		"MONGO_INITDB_ROOT_USERNAME": &m.user,
		"MONGO_INITDB_ROOT_PASSWORD": &m.pass,
		"MONGO_INITDB_DATABASE":      &m.db,
	}

	for ev, v := range envVars {
		var tmp string
		if tmp = os.Getenv(ev); tmp == "" {
			return nil, fmt.Errorf("must set %s in .env", ev)
		}
		*v = tmp
	}
	m.conn = fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin&directConnection=true",
		m.user, m.pass,
		m.host, m.port, m.db,
	)
	return &m, nil
}
