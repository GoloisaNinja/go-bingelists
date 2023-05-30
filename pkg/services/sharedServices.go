package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// General shared func that checks if the media item has already been added to
// the collection name being sent via the func params - favourites/bingelist
// to ensure that duplicate media items do not get added

func AlreadyExists(collectionName, ownerId, mediaId, mediaType string, client *mongo.Client, listId ...primitive.ObjectID) bool {
	var collectionToUse *mongo.Collection
	var filter bson.M
	if collectionName == "favorites" {
		collectionToUse = db.GetCollection(client, "favorites")
		filter = bson.M{"owner": ownerId, "favorites": bson.M{"$elemMatch": bson.M{"mediaId": mediaId, "type": mediaType}}}
	} else {
		collectionToUse = db.GetCollection(client, "bingelists")
		var lid primitive.ObjectID
		if len(listId) > 0 {
			lid = listId[0]
		}
		filter = bson.M{"_id": lid, "titles": bson.M{"$elemMatch": bson.M{"mediaId": mediaId, "type": mediaType}}}
	}
	var result bson.M
	err := collectionToUse.FindOne(context.TODO(), filter).Decode(&result)
	return err == nil
}

// A shared service func that ensures a user being invited to a list, etc - is still
// a valid user that exists in our db

func ValidUser(client *mongo.Client, userId string) bool {
	uc := db.GetCollection(client, "users")
	userIdAsObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false
	}
	filter := bson.M{"_id": userIdAsObj}
	var result bson.M
	err = uc.FindOne(context.TODO(), filter).Decode(&result)
	return err == nil
}

// A shared service func that will be used in both registration middleware as well as
// user change requests to hash a user's submitted password

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// Uses the inbuilt mail package to parse email for validity

func IsValidEmail(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

// A shared service to check if the user email and/or username already exist in the
// db - if either email or name exist - we return true, meaning the user should not
// be recreated or duplicated

func IsExistingUser(name, email string, client *mongo.Client) bool {
	filter := bson.M{"$or": bson.A{bson.M{"email": email}, bson.M{"name": name}}}
	var existingUser models.User
	uc := db.GetCollection(client, "users")
	err := uc.FindOne(context.TODO(), filter).Decode(&existingUser)
	return err == nil
}
