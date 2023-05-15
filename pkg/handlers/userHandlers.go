package handlers

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go-bingelists/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

var secret = util.GetDotEnv("PROD", "JWT_SECRET")

var usersCollection = db.GetCollection(db.DB, "users")
var favoritesCollection = db.GetCollection(db.DB, "favorites")

func FindUserByCredentials(email, password string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return &user, err
	}
	hp := []byte(user.Password)
	p := []byte(password)
	err = bcrypt.CompareHashAndPassword(hp, p)
	return &user, err
}

func GenerateAuthToken(userId string) (string, error) {
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

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	newUser := r.Context().Value("user").(*models.User)
	var resp responses.Response
	_, err := usersCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		resp.Build(500, "internal server error - user not created", nil)
		resp.Respond(w)
		return
	}
	var newUserFavorites models.Favorite
	fId := primitive.NewObjectID()
	newUserFavorites.Build(fId, newUser.Id.Hex())
	_, err = favoritesCollection.InsertOne(context.TODO(), newUserFavorites)
	if err != nil {
		resp.Build(500, "internal server error - favorite creation failed", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "success", newUser)
	resp.Respond(w)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
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
	user, err = FindUserByCredentials(loginReq.Email, loginReq.Password)
	if err != nil {
		resp.Build(403, "email or password invalid", nil)
		resp.Respond(w)
		return
	}
	invalidatedAllTokens := services.InvalidatedAllUserTokens(user.Id.Hex())
	if !invalidatedAllTokens {
		resp.Build(500, "internal server error - token invalidation issue", nil)
		resp.Respond(w)
		return
	}
	var tokenStr string
	tokenStr, err = GenerateAuthToken(user.Id.Hex())
	if err != nil {
		resp.Build(500, "internal server error - token gen failed", nil)
		resp.Respond(w)
		return
	}
	var token models.Token
	tId := primitive.NewObjectID()
	token.Build(tId, tokenStr, user.Id.Hex(), false, false)
	added := services.AddedTokenToCollection(token)
	if !added {
		log.Println("failed to add token to tokensCollection")
	}
	var updatedUser models.User
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": bson.M{"token": token}}
	options := options2.FindOneAndUpdate().SetReturnDocument(options2.After)
	err = usersCollection.FindOneAndUpdate(context.TODO(), filter, update, options).Decode(&updatedUser)
	if err != nil {
		resp.Build(500, "internal server error - problem updating user", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "logged in successfully!", updatedUser)
	resp.Respond(w)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	var resp responses.Response
	userId := r.Context().Value("userId").(string)
	invalidatedTokens := services.InvalidatedAllUserTokens(userId)
	if !invalidatedTokens {
		resp.Build(500, "internal server error - token invalidation failed", nil)
		resp.Respond(w)
		return
	}
	resp.Build(200, "logged out - all tokens expired", nil)
	resp.Respond(w)
}
