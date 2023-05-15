package db

import (
	"context"
	"fmt"
	"go-bingelists/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// prodUri dev envs
//var prodUri = util.GetDotEnv("MONGO_URI")

// devUri dev envs
//var devUri = util.GetDotEnv("MONGO_DEV_URI")

// prodUri prod envs
var prodUri = os.Getenv("MONGO_URI")

// devUri prod envs
var devUri = os.Getenv("MONGO_DEV_URI")
var app config.AppConfig

func ConnectDB() *mongo.Client {
	app.IsProduction = true
	var uriToUse string
	if app.IsProduction {
		uriToUse = prodUri
	} else {
		uriToUse = devUri
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(uriToUse))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("bingelist").Collection(collectionName)
	return coll
}
