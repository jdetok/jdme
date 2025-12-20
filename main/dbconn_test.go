package main

import (
	"fmt"
	"testing"

	"github.com/jdetok/go-api-jdeko.me/api"
)

func TestDBConnection(t *testing.T) {
	app := &api.App{Started: false, QuickStart: true}
	// if err := godotenv.Load("../.env"); err != nil {
	// 	t.Fatalf("fatal error reading .env file: %v", err)
	// }

	if err := app.SetupDB(); err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	fmt.Println("connection to postgres successful")

	// ml, err := logd.NewMongoLogger("log", "http")
	// if err != nil {
	// 	t.Fatalf("failed to connect to mongo: %v", err)
	// }
	// app.Lg.Mongo = ml
	// defer func() {
	// 	if err := app.Lg.Mongo.Client.Disconnect(context.TODO()); err != nil {
	// 		t.Fatalf("fatal mongo error disonnecting: %v", err)
	// 	}
	// }()

}
