package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"social_network_project/api"
	"social_network_project/cache/redisDB"
	"social_network_project/database/postgresql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresqlDB, err := postgresql.ConnectDatabase()
	if err != nil {
		log.Fatal("Error connecting database postgres")
	}

	redisDB, err := redisDB.ConnectToDatabase()
	if err != nil {
		log.Fatal("Error connecting database redis")
	}

	api := api.InitAPI(postgresqlDB, redisDB)

	api.Run(":" + os.Getenv("API_PORT"))
}
