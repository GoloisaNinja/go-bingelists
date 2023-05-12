package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MinifiedBingeList struct {
	ListId string   `bson:"listId" json:"id"`
	Name   string   `bson:"name" json:"name"`
	Movie  []string `bson:"movie" json:"movie"`
	Tv     []string `bson:"tv" json:"tv"`
}

func (mb *MinifiedBingeList) Build(lid, name string) {
	mb.ListId = lid
	mb.Name = name
	mb.Movie = make([]string, 0)
	mb.Tv = make([]string, 0)
}

type BingeList struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Owner      string             `bson:"owner" json:"owner"`
	Users      []string           `bson:"users" json:"users"`
	Titles     []*MediaItem       `bson:"titles" json:"titles"`
	MediaCount int                `bson:"mediaCount" json:"mediaCount"`
	CreatedAt  primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

func (b *BingeList) Build(id primitive.ObjectID, n, o string) {
	b.Id = id
	b.Name = n
	b.Owner = o
	b.Users = make([]string, 0)
	b.Titles = make([]*MediaItem, 0)
	b.MediaCount = 0
	b.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
}
