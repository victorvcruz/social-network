package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"social_network_project/api/request"
)

type Api struct {
	AccountsRequest request.AccountsRequest
}

func (a *Api) Run() {
	router := gin.Default()
	router.POST("/accounts", a.AccountsRequest.CreateAccount)
	router.POST("/accounts/auth", a.AccountsRequest.CreateToken)
	router.GET("/accounts", a.AccountsRequest.GetAccount)
	router.PUT("/accounts", a.AccountsRequest.ChangeAccount)
	router.DELETE("/accounts", a.AccountsRequest.DeleteAccount)

	err := router.Run(":8080")
	if err != nil {
		log.Panic(err)
	}
}
