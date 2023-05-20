package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildMinifiedFavorites(ownerId string, client *mongo.Client) (*models.MinifiedFavorite, error) {
	filter := bson.M{"owner": ownerId}
	var favorites models.Favorite
	fc := db.GetCollection(client, "favorites")
	err := fc.FindOne(context.TODO(), filter).Decode(&favorites)
	if err != nil {
		return nil, err
	}
	var miniFavorite models.MinifiedFavorite
	miniFavorite.Build(ownerId)
	for _, favorite := range favorites.Favorites {
		if favorite.Type == "movie" {
			miniFavorite.Movie = append(miniFavorite.Movie, favorite.MediaId)
		} else {
			miniFavorite.Tv = append(miniFavorite.Tv, favorite.MediaId)
		}
	}
	return &miniFavorite, err
}
