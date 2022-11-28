package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"social_network_project/internal/auth"
	"social_network_project/internal/auth/service"
	"social_network_project/internal/utils/errors"
)

type AuthHandlerClient interface {
	CreateToken(c *gin.Context)
}

type AuthHandler struct {
	service service.AuthServiceClient
}

func RegisterAuthHandler(authService service.AuthServiceClient) AuthHandlerClient {
	return &AuthHandler{
		service: authService,
	}
}
func (a *AuthHandler) CreateToken(c *gin.Context) {
	var request auth.AuthRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	token, err := a.service.CreateToken(request.Email, request.Password)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundEmailError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.UnauthorizedPasswordError:
			log.Println(e)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, token)
	return
}