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
	"io"
	"net/http"
)

func GetFavorites(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		owner := r.Context().Value("userId").(string)
		var favorites models.Favorite
		filter := bson.M{"owner": owner}
		fc := db.GetCollection(c.Config.MongoClient, "favorites")
		err := fc.FindOne(context.TODO(), filter).Decode(&favorites)
		if err != nil {
			resp.Build(500, "internal server error - failed favorite fetch", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success", favorites)
		resp.Respond(w)
	}
}

func GetMinifiedFavorites(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		owner := r.Context().Value("userId").(string)
		miniFavorites, err := services.BuildMinifiedFavorites(owner, c.Config.MongoClient)
		if err != nil {
			resp.Build(500, "internal server error - minified favorite build failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success", miniFavorites)
		resp.Respond(w)
	}
}

func AddToFavorites(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		owner := r.Context().Value("userId").(string)
		var mediaItem models.MediaItem
		rBody, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Build(400, "bad request - media item could not be encoded", nil)
			resp.Respond(w)
			return
		}
		err = json.Unmarshal(rBody, &mediaItem)
		if err != nil {
			resp.Build(400, "bad request - media item could not be encoded", nil)
			resp.Respond(w)
			return
		}
		alreadyExists := services.AlreadyExists("favorites", owner, mediaItem.MediaId, mediaItem.Type, c.Config.MongoClient)
		if alreadyExists {
			resp.Build(400, "bad request - item already in favorites", nil)
			resp.Respond(w)
			return
		}
		var genreName string
		genreName, err = GetGenreNameFromId(mediaItem.Type, c.Config.ApiKey, mediaItem.PrimaryGenreId)
		if err != nil {
			resp.Build(500, "internal server error - problem getting genre from TMDB", nil)
			resp.Respond(w)
			return
		}
		mediaItem.PrimaryGenreName = genreName
		filter := bson.M{"owner": owner}
		update := bson.M{"$push": bson.M{"favorites": mediaItem}}
		fc := db.GetCollection(c.Config.MongoClient, "favorites")
		_, err = fc.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			resp.Build(500, "internal server error - favorite collection not updated", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success - added to favorites", nil)
		resp.Respond(w)
	}
}

func RemoveFromFavorites(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		owner := r.Context().Value("userId").(string)
		mediaId := r.URL.Query().Get("id")
		mediaType := r.URL.Query().Get("type")
		alreadyExists := services.AlreadyExists("favorites", owner, mediaId, mediaType, c.Config.MongoClient)
		if !alreadyExists {
			resp.Build(400, "bad request - item not available to remove", nil)
			resp.Respond(w)
			return
		}
		filter := bson.M{"owner": owner}
		update := bson.M{"$pull": bson.M{"favorites": bson.M{"mediaId": mediaId, "type": mediaType}}}
		fc := db.GetCollection(c.Config.MongoClient, "favorites")
		_, err := fc.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			resp.Build(500, "internal server error - remove favorite failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success - item removed from favorite", nil)
		resp.Respond(w)
	}
}
