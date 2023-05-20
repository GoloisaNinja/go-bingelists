package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type AppConfig struct {
	IsProduction bool
	MongoUri     string
	MongoDevUri  string
	JwtSecret    string
	ApiKey       string
	MongoClient  *mongo.Client
	Port         string
}

func New(isProd bool) *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No env file found...")
	}
	return &AppConfig{
		IsProduction: isProd,
		MongoUri:     getEnv("MONGO_URI", ""),
		MongoDevUri:  getEnv("MONGO_DEV_URI", ""),
		JwtSecret:    getEnv("JWT_SECRET", ""),
		ApiKey:       getEnv("TMDB_APIKEY", ""),
		Port:         ":" + getEnv("PORT", "5000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
