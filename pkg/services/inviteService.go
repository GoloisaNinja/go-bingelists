package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InviteAlreadyAcceptedOrInvitedAlready(client *mongo.Client, listId, invitedId string) bool {
	blc := db.GetCollection(client, "bingelists")
	bingeFilter := bson.M{"_id": listId, "users": invitedId}
	var blResult bson.M
	err := blc.FindOne(context.TODO(), bingeFilter).Decode(&blResult)
	if err == nil {
		return err == nil
	}
	ic := db.GetCollection(client, "invites")
	inviteFilter := bson.M{"bingeListId": listId, "invitedId": invitedId, "pending": true}
	var iResult bson.M
	err = ic.FindOne(context.TODO(), inviteFilter).Decode(&iResult)
	return err == nil
}

func ChangedInvitePendingStatus(client *mongo.Client, inviteId primitive.ObjectID, invitedId string) bool {
	ic := db.GetCollection(client, "invites")
	filter := bson.M{"_id": inviteId, "invitedId": invitedId, "pending": true}
	update := bson.M{"$set": bson.M{"pending": false}}
	result, err := ic.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.ModifiedCount <= 0 {
		return false
	}
	return true
}

func RemovedInviteFromUserInvites(client *mongo.Client, inviteId primitive.ObjectID, userId string) bool {
	uc := db.GetCollection(client, "users")
	userIdAsPrimitive, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false
	}
	filter := bson.M{"_id": userIdAsPrimitive}
	update := bson.M{"$pull": bson.M{"invites": bson.M{"_id": inviteId}}}
	var result *mongo.UpdateResult
	result, err = uc.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.ModifiedCount <= 0 {
		return false
	}
	return true
}

func AcceptedInviteSuccessfully(client *mongo.Client, inviteId primitive.ObjectID, invitedId string) bool {
	ic := db.GetCollection(client, "invites")
	filter := bson.M{"_id": inviteId, "invitedId": invitedId, "pending": true}
	var invite models.Invite
	err := ic.FindOne(context.TODO(), filter).Decode(&invite)
	if err != nil {
		return false
	}
	var blId primitive.ObjectID
	blId, err = primitive.ObjectIDFromHex(invite.BingeListId)
	if err != nil {
		return false
	}
	blc := db.GetCollection(client, "bingelists")
	blFilter := bson.M{"_id": blId}
	update := bson.M{"$push": bson.M{"users": invitedId}}
	var result *mongo.UpdateResult
	result, err = blc.UpdateOne(context.TODO(), blFilter, update)
	if err != nil || result.ModifiedCount <= 0 {
		return false
	}
	return true
}
