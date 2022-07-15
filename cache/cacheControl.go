package cache

import (
	"encoding/json"
	"net/http"
	"os"
	"social_network_project/cache/redisDB"
	"social_network_project/controllers/errors"
	"social_network_project/utils"
)

func InsertCache(req *http.Request, obj any, client *redisDB.RedisClient) error {
	reqID := req.Method + "-" + req.URL.Path + utils.TransformMapInQueryParams(req.URL.Query()) + "-" + req.Header.Get(os.Getenv("JWT_TOKEN_HEADER"))

	respJson, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = client.InsertInDatabase(reqID, string(respJson))
	if err != nil {
		return err
	}

	return nil
}

func FindInCache(req *http.Request, client *redisDB.RedisClient) (*interface{}, error) {
	reqID := req.Method + "-" + req.URL.Path + utils.TransformMapInQueryParams(req.URL.Query()) + "-" + req.Header.Get(os.Getenv("JWT_TOKEN_HEADER"))

	val, err := client.FindInDatabase(reqID)
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
