package handlers

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type JWTCustomClaims struct {
	UserId string
	jwt.RegisteredClaims
}

func FindUserByCredentials(email, password string, client *mongo.Client) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	uc := db.GetCollection(client, "users")
	err := uc.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return &user, err
	}
	hp := []byte(user.Password)
	p := []byte(password)
	err = bcrypt.CompareHashAndPassword(hp, p)
	return &user, err
}

func GenerateAuthToken(userId, secret string) (string, error) {
	jwtSecret := []byte(secret)
	t := time.Now()
	claims := JWTCustomClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(t),
			Issuer:    "bingelists",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, err
}

func CreateNewUser(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newUser := r.Context().Value("user").(*models.User)
		var resp responses.Response
		uc := db.GetCollection(c.Config.MongoClient, "users")
		_, err := uc.InsertOne(context.TODO(), newUser)
		if err != nil {
			resp.Build(500, "internal server error - user not created", nil)
			resp.Respond(w)
			return
		}
		var newUserFavorites models.Favorite
		fId := primitive.NewObjectID()
		newUserFavorites.Build(fId, newUser.Id.Hex())
		fc := db.GetCollection(c.Config.MongoClient, "favorites")
		_, err = fc.InsertOne(context.TODO(), newUserFavorites)
		if err != nil {
			resp.Build(500, "internal server error - favorite creation failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success", newUser)
		resp.Respond(w)
	}
}

func LoginUser(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var loginReq models.LoginRequest
		var resp responses.Response
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			resp.Build(500, "internal server error - user request decode failed", nil)
			resp.Respond(w)
			return
		}
		var user *models.User
		user, err = FindUserByCredentials(loginReq.Email, loginReq.Password, c.Config.MongoClient)
		if err != nil {
			resp.Build(403, "email or password invalid", nil)
			resp.Respond(w)
			return
		}
		invalidatedAllTokens := services.InvalidatedAllUserTokens(user.Id.Hex(), c.Config.MongoClient)
		if !invalidatedAllTokens {
			resp.Build(500, "internal server error - token invalidation issue", nil)
			resp.Respond(w)
			return
		}
		var tokenStr string
		tokenStr, err = GenerateAuthToken(user.Id.Hex(), c.Config.JwtSecret)
		if err != nil {
			resp.Build(500, "internal server error - token gen failed", nil)
			resp.Respond(w)
			return
		}
		var token models.Token
		tId := primitive.NewObjectID()
		token.Build(tId, tokenStr, user.Id.Hex(), false, false)
		added := services.AddedTokenToCollection(token, c.Config.MongoClient)
		if !added {
			log.Println("failed to add token to tokensCollection")
		}
		var updatedUser models.User
		filter := bson.M{"_id": user.Id}
		update := bson.M{"$set": bson.M{"token": token}}
		options := options2.FindOneAndUpdate().SetReturnDocument(options2.After)
		uc := db.GetCollection(c.Config.MongoClient, "users")
		err = uc.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&updatedUser)
		if err != nil {
			resp.Build(500, "internal server error - problem updating user", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "logged in successfully!", updatedUser)
		resp.Respond(w)
	}
}
func Logout(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		userId := r.Context().Value("userId").(string)
		invalidatedTokens := services.InvalidatedAllUserTokens(userId, c.Config.MongoClient)
		if !invalidatedTokens {
			resp.Build(500, "internal server error - token invalidation failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "logged out - all tokens expired", nil)
		resp.Respond(w)
	}
}

func GetPublicUsers(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		userIdStr := r.Context().Value("userId").(string)
		userId, eErr := primitive.ObjectIDFromHex(userIdStr)
		if eErr != nil {
			resp.Build(400, "bad request - invalid user", nil)
			resp.Respond(w)
			return
		}
		uc := db.GetCollection(c.Config.MongoClient, "users")
		filter := bson.M{"isPrivate": false, "_id": bson.M{"$ne": userId}}
		opts := options2.Find().SetProjection(bson.D{{"_id", 1}, {"name", 1}, {"isPrivate", 1}, {"createdAt", 1}})
		cursor, err := uc.Find(context.TODO(), filter, opts)
		defer cursor.Close(context.TODO())
		if err != nil {
			resp.Build(500, "internal server error - user fetch failed", nil)
			resp.Respond(w)
			return
		}
		var users []bson.M
		if err = cursor.All(context.TODO(), &users); err != nil {
			resp.Build(500, "internal server error - user cursor decode failed", nil)
			resp.Respond(w)
			return
		}
		resp.Build(200, "success", users)
		resp.Respond(w)
	}
}

func UserNameChange(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type nameRequest struct {
			Name string `json:"name"`
		}
		var resp responses.Response
		userIdStr := r.Context().Value("userId").(string)
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			resp.Build(500, "internal server error - problem converting user to primitive", nil)
			resp.Respond(w)
			return
		}
		var newName nameRequest
		err = json.NewDecoder(r.Body).Decode(&newName)
		if err != nil {
			resp.Build(400, "bad request - json decoding error", nil)
			resp.Respond(w)
			return
		}
		if len(newName.Name) > 10 {
			resp.Build(400, "bad request - name exceeds max length", nil)
			resp.Respond(w)
			return
		}
		var user models.User
		uc := db.GetCollection(c.Config.MongoClient, "users")
		filter := bson.M{"_id": userId}
		err = uc.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			resp.Build(400, "bad request - invalid user", nil)
			resp.Respond(w)
			return
		}
		// has the user submitted a change request with no actual change to name?
		if user.FirstName == newName.Name {
			resp.Build(400, "bad request - names are equal", nil)
			resp.Respond(w)
			return
		}
		// check if the user name submitted already exists
		duplicatesFilter := bson.M{"name": newName.Name}
		var existingUserWithSameName bson.M
		err = uc.FindOne(context.TODO(), duplicatesFilter).Decode(&existingUserWithSameName)
		if err == nil {
			resp.Build(400, "bad request - user name already taken", nil)
			resp.Respond(w)
			return
		}
		// if we make here - all error checks for duplicates and invalid changes have passed
		update := bson.M{"$set": bson.M{"name": newName.Name}}
		updateSuccess, err := uc.UpdateOne(context.TODO(), filter, update)
		if err != nil || updateSuccess.ModifiedCount <= 0 {
			resp.Build(500, "internal server error - name change failed", nil)
			resp.Respond(w)
			return
		}
		user.FirstName = newName.Name
		resp.Build(200, "success", user)
		resp.Respond(w)
	}
}

func UserPrivacyChange(c *config.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp responses.Response
		userIdStr := r.Context().Value("userId").(string)
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			resp.Build(400, "bad request - invalid user", nil)
			resp.Respond(w)
			return
		}
		type privacyRequest struct {
			IsPrivate bool `json:"isPrivate"`
		}
		var privacyChange privacyRequest
		err = json.NewDecoder(r.Body).Decode(&privacyChange)
		if err != nil {
			resp.Build(400, "bad request - error decoding privacy", nil)
			resp.Respond(w)
			return
		}
		var user models.User
		filter := bson.M{"_id": userId}
		uc := db.GetCollection(c.Config.MongoClient, "users")
		err = uc.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			resp.Build(400, "bad request - invalid user", nil)
			resp.Respond(w)
			return
		}
		if privacyChange.IsPrivate == user.IsPrivate {
			resp.Build(400, "bad request - privacy is equal", nil)
			resp.Respond(w)
			return
		}
		update := bson.M{"$set": bson.M{"isPrivate": privacyChange.IsPrivate}}
		var updateSuccess *mongo.UpdateResult
		updateSuccess, err = uc.UpdateOne(context.TODO(), filter, update)
		if err != nil || updateSuccess.ModifiedCount <= 0 {
			resp.Build(500, "internal server error - privacy update failed", nil)
			resp.Respond(w)
			return
		}
		user.IsPrivate = privacyChange.IsPrivate
		resp.Build(200, "success", user)
		resp.Respond(w)
	}
}
