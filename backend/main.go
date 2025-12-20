package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/mukesh1352/splitwise-backend/db"
)

func main() {
	// Load .env file (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	database := db.Connect()
	defer database.Close()

	log.Println("application started successfully")
}

