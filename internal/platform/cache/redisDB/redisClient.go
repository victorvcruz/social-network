package redisDB

import (
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"social_network_project/internal/utils/errors"
	"time"
)


type RedisClient interface {
	ConnectToDatabase() error
	InsertInDatabase(key string, value string) error
	FindInDatabase(key string) (string, error)
}

type Redis struct {
	Client *redis.Client
}

func NewRedis() RedisClient {
	return &Redis{}
}

func (r *Redis) ConnectToDatabase() error {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return err
	}
	log.Println("Redis Connected")
	r.Client = client
	return nil
}

func (r *Redis) InsertInDatabase(key string, value string) error {
	err := r.Client.Set(r.Client.Context(), key, value, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) FindInDatabase(key string) (string, error) {
	val, err := r.Client.Get(r.Client.Context(), key).Result()
	if err != nil {
		return "", &errors.CacheNotFoundError{}
	}
	return val, nil
}
