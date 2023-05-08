package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetDotEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in env load")
	}
	return os.Getenv(key)
}
