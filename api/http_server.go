package api

import (
	"github.com/gin-gonic/gin"
	"social_network_project/api/handler"
	"social_network_project/controllers"
)

func InitAPI() *gin.Engine {
	ginServer := gin.Default()

	handler.RegisterAccountsHandlers(ginServer, controllers.NewAccountsController())
	handler.RegisterPostsHandlers(ginServer, controllers.NewPostsController())

	return ginServer
}
