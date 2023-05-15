package db

import (
	"context"
	"fmt"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var prodUri = util.GetDotEnv("MONGO_URI")
var devUri = util.GetDotEnv("MONGO_DEV_URI")
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
