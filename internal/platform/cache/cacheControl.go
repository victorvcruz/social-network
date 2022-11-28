package cache

import (
	"encoding/json"
	"net/http"
	"os"
	"social_network_project/internal/platform/cache/redisDB"
	"social_network_project/internal/utils"
	"social_network_project/internal/utils/errors"
)

type RedisServiceClient interface {
	InsertCache(req *http.Request, obj any) error
	FindInCache(req *http.Request) (*interface{}, error)
}

func NewRedisService(_client redisDB.RedisClient) RedisServiceClient {
	return &RedisService{
		client: _client,
	}
}

type RedisService struct {
	client redisDB.RedisClient
}

func (r *RedisService) InsertCache(req *http.Request, obj any) error {
	reqID := req.Method + "-" + req.URL.Path + utils.TransformMapInQueryParams(req.URL.Query()) + "-" + req.Header.Get(os.Getenv("JWT_TOKEN_HEADER"))

	respJson, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = r.client.InsertInDatabase(reqID, string(respJson))
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisService) FindInCache(req *http.Request) (*interface{}, error) {
	reqID := req.Method + "-" + req.URL.Path + utils.TransformMapInQueryParams(req.URL.Query()) + "-" + req.Header.Get(os.Getenv("JWT_TOKEN_HEADER"))

	val, err := r.client.FindInDatabase(reqID)
	if err != nil {
		switch e := err.(type) {
		case *errors.CacheNotFoundError:
			return nil, e
		}
	}

	var responseCache interface{}

	err = json.Unmarshal([]byte(val), &responseCache)
	if err != nil {
		return nil, err
	}

	return &responseCache, nil
}
