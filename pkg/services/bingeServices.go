package services

import (
	"context"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildMinifiedBingeSlicesByOwner(owner string, client *mongo.Client) ([]*models.MinifiedBingeList, error) {
	filter := bson.M{"owner": owner}
	bc := db.GetCollection(client, "bingelists")
	cursor, err := bc.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var lists []*models.BingeList
	var miniLists []*models.MinifiedBingeList
	if err = cursor.All(context.TODO(), &lists); err != nil {
		return nil, err
	}
	for _, list := range lists {
		var tempList models.MinifiedBingeList
		tempList.Build(list.Id.Hex(), list.Name)
		for _, title := range list.Titles {
			if title.Type == "movie" {
				tempList.Movie = append(tempList.Movie, title.MediaId)
			} else {
				tempList.Tv = append(tempList.Tv, title.MediaId)
			}
		}
		miniLists = append(miniLists, &tempList)
	}
	return miniLists, err
}
