package main

import (
	"fmt"
	"go-bingelists/pkg/config"
	"go-bingelists/pkg/db"
	"log"
	"net/http"
	"os"
)

var app config.AppConfig

func main() {
	app.IsProduction = false
	db.ConnectDB()
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":5000"
	}
	srv := &http.Server{
		Addr:    port,
		Handler: routes(),
	}
	fmt.Println("Server is up on port " + port)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
