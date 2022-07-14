package redis

import (
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"social_network_project/controllers/errors"
	"time"
)

var client *redis.Client

func ConnectToDatabase() error {

	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return err
	}

	log.Println("Redis Connected")
	return nil
}

func InsertInDatabase(key string, value string) error {
	err := client.Set(client.Context(), key, value, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func FindInDatabase(key string) (string, error) {
	val, err := client.Get(client.Context(), key).Result()
	if err != nil {
		return "", &errors.CacheNotFoundError{}
	}
	return val, nil
}
