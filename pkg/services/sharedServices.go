package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AlreadyExists(collectionName, ownerId, mediaId, mediaType string, listId ...primitive.ObjectID) bool {
	var collectionToUse *mongo.Collection
	var filter bson.M
	if collectionName == "favorites" {
		collectionToUse = favoritesCollection
		filter = bson.M{"owner": ownerId, "favorites": bson.M{"$elemMatch": bson.M{"mediaId": mediaId, "type": mediaType}}}
	} else {
		collectionToUse = bingeListCollection
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