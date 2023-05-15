package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-bingelists/pkg/db"
	"go-bingelists/pkg/models"
	"go-bingelists/pkg/responses"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"regexp"
)

var tokensCollection = db.GetCollection(db.DB, "tokens")

var secret = os.Getenv("JWT_SECRET")

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var resp responses.Response
			x, err := regexp.Compile(`^(?P<B>Bearer\s+)(?P<T>.*)$`)
			if err != nil {
				resp.Build(403, "unauthorized - no token", nil)
				resp.Respond(w)
				return
			}
			authHeader := x.FindStringSubmatch(r.Header.Get("Authorization"))
			if len(authHeader) != 3 {
				resp.Build(403, "unauthorized - bad header", nil)
				resp.Respond(w)
				return
			}
			tokenIndex := x.SubexpIndex("T")
			tokenString := authHeader[tokenIndex]
			var token *jwt.Token
			token, err = jwt.Parse(
				tokenString,
				func(token *jwt.Token) (interface{}, error) {
					_, ok := token.Method.(*jwt.SigningMethodHMAC)
					if !ok {
						authErr := errors.New("unauthorized")
						return nil, authErr
					}
					return []byte(secret), nil
				},
			)
			if err != nil {
				resp.Build(403, err.Error(), nil)
				resp.Respond(w)
				return
			}
			var userId string
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok && token.Valid {
				userId = claims["UserId"].(string)
			}
			filter := bson.M{"token": tokenString}
			var dbToken models.Token
			err = tokensCollection.FindOne(context.TODO(), filter).Decode(&dbToken)
			if err != nil {
				resp.Build(403, "unauthorized token", nil)
				resp.Respond(w)
				return
			}
			if dbToken.IsRevoked || dbToken.IsExpired {
				resp.Build(403, "unauthorized token", nil)
				resp.Respond(w)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
