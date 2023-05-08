package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

var tokenCollection = db.GetCollection(db.DB, "tokens")

func AddedTokenToCollection(t models.Token) bool {
	_, err := tokenCollection.InsertOne(context.TODO(), t)
	return err == nil
}

func InvalidatedAllOtherUserTokensOnLogin(userId string) bool {
	filter := bson.M{"user": userId, "isRevoked": false, "isExpired": false}
	update := bson.M{"$set": bson.M{"isExpired": true, "isRevoked": true}}
	_, err := tokenCollection.UpdateMany(context.TODO(), filter, update)
	return err == nil
}
