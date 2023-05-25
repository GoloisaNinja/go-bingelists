package handlers

import (
	"context"
	"encoding/json"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func GetPendingInvites(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		userId := r.Context().Value("userId").(string)
		ic := db.GetCollection(c.Config.MongoClient, "invites")
		filter := bson.M{"invitedById": userId, "pending": true}
		cursor, err := ic.Find(context.TODO(), filter)
		if err != nil {
			resp.Build(500, "internal server error - invite fetch failed", nil)
			resp.Respond(w)
			return
		}
		var invites []models.Invite
		err = cursor.All(context.TODO(), &invites)
		if err != nil {
			resp.Build(500, "internal server error - invite decode failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success", invites)
		resp.Respond(w)
	}
}

func CreateNewInvite(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		var invite models.Invite
		err := json.NewDecoder(r.Body).Decode(&invite)
		if err != nil {
			resp.Build(400, "bad request - invalid invite", nil)
			resp.Respond(w)
			return
		}
		// check if invite already exists...
		duplicate := services.InviteAlreadyAcceptedOrInvitedAlready(c.Config.MongoClient, invite.BingeListId, invite.InvitedId)
		if duplicate {
			resp.Build(400, "invalid invite - invite pending or user already on list", nil)
			resp.Respond(w)
			return
		}
		// check that the invited user is valid
		validUser := services.ValidUser(c.Config.MongoClient, invite.InvitedId)
		if !validUser {
			resp.Build(400, "bad request - invalid user invited", nil)
			resp.Respond(w)
			return
		}
		// if not finish building valid invite...
		invite.Id = primitive.NewObjectID()
		invite.Pending = true
		ic := db.GetCollection(c.Config.MongoClient, "invites")
		_, err = ic.InsertOne(context.TODO(), invite)
		if err != nil {
			resp.Build(500, "internal server error - invite insert failed", nil)
			resp.Respond(w)
			return
		}
		uc := db.GetCollection(c.Config.MongoClient, "users")
		var uidAsObj primitive.ObjectID
		uidAsObj, err = primitive.ObjectIDFromHex(invite.InvitedId)
		if err != nil {
			resp.Build(500, "internal server error - id encode failure", nil)
			resp.Respond(w)
			return
		}
		filter := bson.M{"_id": uidAsObj}
		update := bson.M{"$push": bson.M{"invites": invite}}
		var result *mongo.UpdateResult
		result, err = uc.UpdateOne(context.TODO(), filter, update)
		if err != nil || result.ModifiedCount <= 0 {
			resp.Build(500, "internal server error - invite to user failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success - invite created and sent", nil)
		resp.Respond(w)
	}
}

func ProcessInvite(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		userIdAsStr := r.Context().Value("userId").(string)
		inviteIdAsStr := r.URL.Query().Get("id")
		action := r.URL.Query().Get("action")
		inviteId, err := primitive.ObjectIDFromHex(inviteIdAsStr)
		if err != nil {
			resp.Build(400, "bad request - invalid invite", nil)
			resp.Respond(w)
			return
		}
		if action == "accept" {
			acceptSuccess := services.AcceptedInviteSuccessfully(c.Config.MongoClient, inviteId, userIdAsStr)
			if !acceptSuccess {
				resp.Build(500, "internal server error - problem accepting invite as user", nil)
				resp.Respond(w)
				return
			}
		}
		changedPending := services.ChangedInvitePendingStatus(c.Config.MongoClient, inviteId, userIdAsStr)
		if !changedPending {
			resp.Build(500, "internal server error - problem with pending status", nil)
			resp.Respond(w)
			return
		}
		removedInvite := services.RemovedInviteFromUserInvites(c.Config.MongoClient, inviteId, userIdAsStr)
		if !removedInvite {
			resp.Build(500, "internal server error - problem with user invites object", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success - processed invite", nil)
		resp.Respond(w)
	}
}
