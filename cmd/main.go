package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"go-bingelists/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	db.ConnectDB()
	port := ":" + os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	}
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization", "Bearer", "Accept", "Accept-Language", "Origin", "Accept-Encoding", "Content-Length", "Referrer", "User-Agent"})
	originOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "DELETE", "OPTIONS"})
	srv := &http.Server{
		Addr:    port,
		Handler: handlers.CORS(originOk, headersOk, methodsOk)(routes()),
	}
	fmt.Println("Server is up on port " + port)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
