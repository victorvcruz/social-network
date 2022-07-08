package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"social_network_project/api/handler"
)

type Api struct {
	AccountsAPI request.AccountsAPI
}

func (a *Api) Run() {
	router := gin.Default()
	router.POST("/accounts", a.AccountsAPI.CreateAccount)
	router.POST("/accounts/auth", a.AccountsAPI.CreateToken)
	router.GET("/accounts", a.AccountsAPI.GetAccount)
	router.PUT("/accounts", a.AccountsAPI.UpdateAccount)
	router.DELETE("/accounts", a.AccountsAPI.DeleteAccount)

	err := router.Run(":8080")
	if err != nil {
		log.Panic(err)
	}
}
