package db

import (
	"context"
	"fmt"
	"go-bingelists/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDB(r *config.Repository) *mongo.Client {
	var uriToUse string
	if r.Config.IsProduction {
		uriToUse = r.Config.MongoUri
	} else {
		uriToUse = r.Config.MongoDevUri
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

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("bingelist").Collection(collectionName)
	return coll
}
