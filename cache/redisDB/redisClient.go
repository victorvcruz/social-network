package redisDB

import (
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"social_network_project/controllers/errors"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func ConnectToDatabase() (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis Connected")
	return client, nil
}

func (r *RedisClient) InsertInDatabase(key string, value string) error {
	err := r.Client.Set(r.Client.Context(), key, value, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) FindInDatabase(key string) (string, error) {
	val, err := r.Client.Get(r.Client.Context(), key).Result()
	if err != nil {
		return "", &errors.CacheNotFoundError{}
	}
	return val, nil
}
