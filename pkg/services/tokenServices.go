package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddedTokenToCollection(t models.Token, client *mongo.Client) bool {
	tc := db.GetCollection(client, "tokens")
	_, err := tc.InsertOne(context.TODO(), t)
	return err == nil
}

func InvalidatedAllUserTokens(userId string, client *mongo.Client) bool {
	filter := bson.M{"user": userId, "isRevoked": false, "isExpired": false}
	update := bson.M{"$set": bson.M{"isExpired": true, "isRevoked": true}}
	tc := db.GetCollection(client, "tokens")
	_, err := tc.UpdateMany(context.TODO(), filter, update)
	return err == nil
}
