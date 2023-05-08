package middleware

import (
	"context"
	"encoding/json"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/handlers"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
)

var usersCollection = db.GetCollection(db.DB, "users")

func isValidEmail(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

func isExistingUser(email string) bool {
	filter := bson.M{"email": email}
	var existingUser models.User
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&existingUser)
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func Registration(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var resp responses.Response
			var newReq models.NewUserRequest
			defer r.Body.Close()
			err := json.NewDecoder(r.Body).Decode(&newReq)
			if err != nil {
				resp.Build(400, "bad request", nil)
				resp.Respond(w)
				return
			}
			if newReq.Email == "" || newReq.Password == "" {
				resp.Build(400, "bad request", nil)
				resp.Respond(w)
				return
			}
			if len(newReq.Password) < 9 {
				resp.Build(400, "password length insufficient", nil)
				resp.Respond(w)
				return
			}
			validEmail := isValidEmail(newReq.Email)
			exists := isExistingUser(newReq.Email)
			if !validEmail || exists {
				resp.Build(400, "user or email invalid", nil)
				resp.Respond(w)
				return
			}
			hashed, err := hashPassword(newReq.Password)
			if err != nil {
				resp.Build(500, "internal server error", nil)
				resp.Respond(w)
				return
			}
			var newUser models.User
			uid := primitive.NewObjectID()
			tokenStr, err := handlers.GenerateAuthToken(uid.Hex())
			if err != nil {
				resp.Build(500, "internal server error - token creation failed", nil)
				resp.Respond(w)
				return
			}
			tid := primitive.NewObjectID()
			var token models.Token
			token.Build(tid, tokenStr, uid.Hex(), false, false)
			addedToken := services.AddedTokenToCollection(token)
			if !addedToken {
				log.Println("error writing token to db...")
			}
			newUser.Build(uid, newReq.Name, newReq.Email, hashed, newReq.IsPrivate, token)
			ctx := context.WithValue(r.Context(), "user", &newUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
