package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Token struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Token     string             `bson:"token" json:"token"`
	IsRevoked bool               `bson:"isRevoked" json:"isRevoked"`
	IsExpired bool               `bson:"isExpired" json:"isExpired"`
	User      string             `bson:"user" json:"user"`
}

func (t *Token) Build(id primitive.ObjectID, tStr, u string, r, e bool) {
	t.Id = id
	t.Token = tStr
	t.IsRevoked = r
	t.IsExpired = e
	t.User = u
}

// TODO - add binglelist criteria to invite after bingelist model is created

type Invite struct {
	Id primitive.ObjectID
}

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	IsPrivate bool               `bson:"isPrivate" json:"isPrivate"`
	Token     *Token             `bson:"token" json:"token"`
	Invites   []*Invite          `bson:"invites" json:"invites"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

func (u *User) Build(id primitive.ObjectID, n, e, hp string, ip bool, token Token) {
	u.Id = id
	u.FirstName = n
	u.Email = e
	u.Password = hp
	u.IsPrivate = ip
	u.Token = &token
	u.Invites = make([]*Invite, 0)
	u.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
}

type NewUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsPrivate bool   `json:"isPrivate"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
