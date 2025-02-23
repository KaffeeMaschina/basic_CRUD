package main

import (
	"github.com/KaffeeMaschina/basic_CRUD/internal/app/apiserver"
	"github.com/KaffeeMaschina/basic_CRUD/internal/app/storage"
	"log"
)

func main() {
	// Loading config from .env
	config := apiserver.LoadConfig()
	// Get connection to Database
	db, err := storage.NewDatabase(config.DBUser, config.DBPass, config.DBPort, config.DBName)
	if err != nil {
		log.Fatal(err)
	}
	// Create new server and start
	s := apiserver.New(config, db)
	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
}
