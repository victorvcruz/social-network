package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"social_network_project/api/handler"
	"social_network_project/cache/redisDB"
	"social_network_project/controllers"
	message_broker "social_network_project/message-broker"
)

func InitAPI(postgresDB *sql.DB, redis *redis.Client, rabbitConn *amqp.Connection) *gin.Engine {
	ginServer := gin.Default()

	redisClient := &redisDB.RedisClient{redis}
	rabbitmq := &message_broker.NotificationControl{
		rabbitConn,
		message_broker.NewNotificationController(postgresDB),
	}

	go rabbitmq.ConsumerMessage()

	handler.RegisterAccountsHandlers(ginServer, controllers.NewAccountsController(postgresDB, rabbitmq), redisClient)
	handler.RegisterPostsHandlers(ginServer, controllers.NewPostsController(postgresDB, rabbitmq), redisClient)
	handler.RegisterCommentsHandlers(ginServer, controllers.NewCommentsController(postgresDB, rabbitmq), redisClient)
	handler.RegisterInteractionsHandlers(ginServer, controllers.NewInteractionsController(postgresDB, rabbitmq))
	return ginServer
}
