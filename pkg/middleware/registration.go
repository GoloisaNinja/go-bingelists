package middleware

import (
	"context"
	"encoding/json"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/handlers"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go-bingelists/pkg/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func Registration(c *config.Repository, next http.Handler) http.Handler {
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
			validEmail := services.IsValidEmail(newReq.Email)
			if !validEmail {
				resp.Build(400, "email invalid", nil)
				resp.Respond(w)
				return
			}
			exists := services.IsExistingUser(newReq.Name, newReq.Email, c.Config.MongoClient)
			if exists {
				resp.Build(400, "username or email already exists", nil)
				resp.Respond(w)
				return
			}
			hashed, err := services.HashPassword(newReq.Password)
			if err != nil {
				resp.Build(500, "internal server error", nil)
				resp.Respond(w)
				return
			}
			var newUser models.User
			uid := primitive.NewObjectID()
			tokenStr, err := handlers.GenerateAuthToken(uid.Hex(), c.Config.JwtSecret)
			if err != nil {
				resp.Build(500, "internal server error - token creation failed", nil)
				resp.Respond(w)
				return
			}
			tid := primitive.NewObjectID()
			var token models.Token
			token.Build(tid, tokenStr, uid.Hex(), false, false)
			addedToken := services.AddedTokenToCollection(token, c.Config.MongoClient)
			if !addedToken {
				log.Println("error writing token to db...")
			}
			newUser.Build(uid, newReq.Name, newReq.Email, hashed, newReq.IsPrivate, token)
			ctx := context.WithValue(r.Context(), "user", &newUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
