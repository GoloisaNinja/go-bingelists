package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Favorite struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Owner     string             `bson:"owner" json:"owner"`
	Favorites []*MediaItem       `bson:"favorites" json:"favorites"`
}

func (f *Favorite) Build(id primitive.ObjectID, o string) {
	f.Id = id
	f.Owner = o
	f.Favorites = make([]*MediaItem, 0)
}

type MinifiedFavorite struct {
	Owner string   `json:"owner"`
	Movie []string `json:"movie"`
	Tv    []string `json:"tv"`
}

func (mf *MinifiedFavorite) Build(o string) {
	mf.Owner = o
	mf.Movie = make([]string, 0)
	mf.Tv = make([]string, 0)
}
