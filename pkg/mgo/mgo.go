package mgo

import (
	"context"
	"fmt"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/conn"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin&directConnection=true",
		e.User, e.Pass, e.Host, e.Port, e.Database), nil
}
