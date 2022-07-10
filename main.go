package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"social_network_project/api"
	"social_network_project/database/postgresql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresql.ConnectDatabase()

	api := api.InitAPI()

	api.Run(":" + os.Getenv("API_PORT"))
}
