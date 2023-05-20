package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/db"
	"log"
	"net/http"
)

func main() {
	appConfig := config.New(true)

	fmt.Println("app Production is: ", appConfig.IsProduction)
	// Create config repo to be shared through application
	configRepo := config.NewRepo(appConfig)
	config.NewAppConfiguration(configRepo)
	// Set mongo client in app config
	appConfig.MongoClient = db.ConnectDB(configRepo)
	// dev database seed of fake users for edge casing certain FE components
	//data.SeedDB()

	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization", "Bearer", "Accept", "Accept-Language", "Origin", "Accept-Encoding", "Content-Length", "Referrer", "User-Agent"})
	originOk := handlers.AllowedOrigins([]string{"http://localhost:3000", "https://bingelists.netlify.app", "https://bingelists.app"})
	methodsOk := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "DELETE", "OPTIONS"})
	srv := &http.Server{
		Addr:    appConfig.Port,
		Handler: handlers.CORS(originOk, headersOk, methodsOk)(routes(configRepo)),
	}
	fmt.Println("Server is up on port " + appConfig.Port)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
