package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	dbClient, err := database.NewDatabaseClient(cfg)
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}

	defer dbClient.DB.Close()
}
