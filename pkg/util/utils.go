package util

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-bingelists/pkg/config"
	"log"
	"os"
)

var app config.AppConfig

func SetUtilConfig(a *config.AppConfig) {
	app = *a
}

func GetDotEnv(key string) string {
	isProd := true
	if !isProd {
		fmt.Println(app.IsProduction)
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error in env load")
		}
	}
	return os.Getenv(key)
}
