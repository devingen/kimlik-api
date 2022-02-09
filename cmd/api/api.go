package main

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/kimlik-api/config"
	"github.com/devingen/kimlik-api/server"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

func main() {

	var appConfig config.App
	err := envconfig.Process("kimlik_api", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Running server on port", appConfig.Port)

	db, err := database.New(appConfig.Mongo.URI)
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	srv := server.New(appConfig, db)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and serve failed %s", err.Error())
	}
}
