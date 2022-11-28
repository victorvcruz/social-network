package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"social_network_project/cmd/api"
	"social_network_project/cmd/api/handlers"
	"social_network_project/internal/account"
	service9 "social_network_project/internal/account/service"
	service4 "social_network_project/internal/auth/service"
	service3 "social_network_project/internal/comment"
	"social_network_project/internal/comment/service"
	service2 "social_network_project/internal/interaction"
	service6 "social_network_project/internal/interaction/service"
	"social_network_project/internal/notification"
	service7 "social_network_project/internal/notification/service"
	"social_network_project/internal/platform/cache"
	"social_network_project/internal/platform/cache/redisDB"
	"social_network_project/internal/platform/database/postgresql"
	"social_network_project/internal/platform/message-broker/rabbitmq"
	service5 "social_network_project/internal/post"
	service8 "social_network_project/internal/post/service"
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

	redisDB := redisDB.NewRedis()
	err = redisDB.ConnectToDatabase()
	if err != nil {
		log.Fatal("Error connecting database redis")
	}

	redisService := cache.NewRedisService(redisDB)

	rabbitConn, err := rabbitmq.ConnectToMessageBroker()
	if err != nil {
		log.Fatal("Error connecting message-broker rabbitMQ")
	}
	notificationRepository := notification.NewNotificationRepository(postgresqlDB)
	notificationService := service7.NewNotificationService(rabbitConn, notificationRepository)
	go notificationService.ConsumerMessage()

	accountsRepository := account.NewAccountRepository(postgresqlDB)
	postsRepository :=    service5.NewPostRepository(postgresqlDB)
	commentsRepository := service3.NewComentRepository(postgresqlDB)
	interactionsRepository := service2.NewInteractionRepository(postgresqlDB)

	authService := service4.NewAuthService(accountsRepository)
	accountsService := service9.NewAccountsService(accountsRepository, notificationService)
	postsService := service8.NewPostsService(postsRepository, accountsRepository, notificationService)
	commentsService := service.NewCommentsService(commentsRepository, accountsRepository, postsRepository, notificationService)
	interactionsService := service6.NewInteractionsService(accountsRepository, commentsRepository, interactionsRepository, notificationService)

	authHandler := handlers.RegisterAuthHandler(authService)
	accountsHandler := handlers.RegisterAccountsHandlers(accountsService, redisService)
	postsHandler := handlers.RegisterPostsHandlers(postsService, redisService)
	commentsHandler := handlers.RegisterCommentsHandlers(commentsService, redisService)
	interactionsHandler := handlers.RegisterInteractionsHandlers(interactionsService)

	api := api.Init(authHandler, accountsHandler, postsHandler, commentsHandler, interactionsHandler)
	api.Run(":" + os.Getenv("API_PORT"))
}
