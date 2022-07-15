package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"social_network_project/api/handler"
	"social_network_project/cache/redisDB"
	"social_network_project/controllers"
)

func InitAPI(postgresDB *sql.DB, redis *redis.Client) *gin.Engine {
	ginServer := gin.Default()

	redisClient := &redisDB.RedisClient{redis}

	handler.RegisterAccountsHandlers(ginServer, controllers.NewAccountsController(postgresDB), redisClient)
	handler.RegisterPostsHandlers(ginServer, controllers.NewPostsController(postgresDB), redisClient)
	handler.RegisterCommentsHandlers(ginServer, controllers.NewCommentsController(postgresDB), redisClient)
	handler.RegisterInteractionsHandlers(ginServer, controllers.NewInteractionsController(postgresDB))

	return ginServer
}
