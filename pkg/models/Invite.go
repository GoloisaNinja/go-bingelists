package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Invite struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	BingeListId   string             `bson:"bingeListId" json:"bingeListId"`
	BingeListName string             `bson:"bingeListName" json:"bingeListName"`
	InvitedById   string             `bson:"invitedById" json:"invitedById"`
	InvitedByName string             `bson:"invitedByName" json:"invitedByName"`
	InvitedId     string             `bson:"invitedId" json:"invitedId"`
	InvitedName   string             `bson:"invitedName" json:"invitedName"`
	Message       string             `bson:"message" json:"message"`
	Pending       bool               `bson:"pending" json:"pending"`
}
