package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	account2 "social_network_project/internal/account"
	"social_network_project/internal/account/service"
	"social_network_project/internal/platform/cache"
	"social_network_project/internal/utils"
	"social_network_project/internal/utils/errors"
	"social_network_project/internal/utils/validate"
	"strconv"
	"time"
)

type AccountsHandlerClient interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	UpdateAccount(c *gin.Context)
	DeleteAccount(c *gin.Context)
	FollowAccount(c *gin.Context)
	SearchFollowing(c *gin.Context)
	SearchFollowers(c *gin.Context)
	UnfollowAccount(c *gin.Context)
}

type AccountsHandler struct {
	Controller  service.AccountsServiceClient
	RedisClient cache.RedisServiceClient
	Validate    *validator.Validate
}

func RegisterAccountsHandlers(accountsController service.AccountsServiceClient, _redis cache.RedisServiceClient) AccountsHandlerClient {
	return &AccountsHandler{
		Controller:  accountsController,
		RedisClient: _redis,
		Validate:    validator.New(),
	}
}

func (a *AccountsHandler) CreateAccount(c *gin.Context) {
	var request account2.AccountRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	account := a.fillFields(request)
	mapper := make(map[string]interface{})
	err = a.Validate.Struct(account)
	if err != nil {
		mapper["errors"] = validate.RequestAccountValidate(err)
		c.JSON(http.StatusBadRequest, mapper)

		return
	}

	err = a.Controller.InsertAccount(account)
	if err != nil {
		switch e := err.(type) {
		case *errors.ConflictUsernameError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return

		case *errors.ConflictEmailError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsHandler) GetAccount(c *gin.Context) {

	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	resCache, err := a.RedisClient.FindInCache(c.Request)
	switch e := err.(type) {
	case *errors.CacheNotFoundError:
		log.Println(e.Error())
	default:
		c.JSON(http.StatusOK, resCache)
		log.Println("Cache")

		return
	}

	account, err := a.Controller.FindAccountByID(&id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, account.ToResponse())

	c.JSON(http.StatusOK, account.ToResponse())
	return

}

func (a *AccountsHandler) UpdateAccount(c *gin.Context) {

	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request account2.AccountRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	account, err := a.Controller.FindAccountByID(&id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	accountChange := a.mergeAccountToUpdatedAccount(account, request)

	mapper := make(map[string]interface{})

	err = a.Validate.Struct(accountChange)
	if err != nil {
		mapper["errors"] = validate.RequestAccountValidate(err)
		c.JSON(http.StatusBadRequest, mapper)

		return
	}

	err = a.Controller.ChangeAccountDataByID(&id, request)
	if err != nil {
		switch e := err.(type) {
		case *errors.ConflictUsernameError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return

		case *errors.ConflictEmailError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsHandler) DeleteAccount(c *gin.Context) {
	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	account, err := a.Controller.DeleteAccountByID(&id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsHandler) FollowAccount(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request account2.AccountRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	if request.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Add ID",
		})
		return
	}

	accountToFollow := request.ID

	account, err := a.Controller.CreateFollow(&accountID, &accountToFollow)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.ConflictAlreadyFollowError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsHandler) SearchFollowing(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	resCache, err := a.RedisClient.FindInCache(c.Request)
	switch e := err.(type) {
	case *errors.CacheNotFoundError:
		log.Println(e.Error())
	default:
		c.JSON(http.StatusOK, resCache)
		log.Println("Cache")

		return
	}

	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Page is not a number",
		})
		return
	}
	listOfAccounts, err := a.Controller.FindAccountsFollowing(&accountID, &page)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, listOfAccounts)

	c.JSON(http.StatusOK, listOfAccounts)
	return
}

func (a *AccountsHandler) SearchFollowers(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	resCache, err := a.RedisClient.FindInCache(c.Request)
	switch e := err.(type) {
	case *errors.CacheNotFoundError:
		log.Println(e.Error())
	default:
		c.JSON(http.StatusOK, resCache)
		log.Println("Cache")

		return
	}

	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Page is not a number",
		})
		return
	}
	listOfAccounts, err := a.Controller.FindAccountsFollowers(&accountID, &page)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, listOfAccounts)

	c.JSON(http.StatusOK, listOfAccounts)
	return
}

func (a *AccountsHandler) UnfollowAccount(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request account2.AccountRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	accountToFollow := request.ID

	account, err := a.Controller.DeleteFollow(&accountID, &accountToFollow)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.ConflictAlreadyUnfollowError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsHandler) mergeAccountToUpdatedAccount(account *account2.Account, req account2.AccountRequest) *account2.Account {

	if req.Username != "" {
		account.Username = req.Username
	}
	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Description != "" {
		account.Description = req.Description
	}
	if req.Email != "" {
		account.Email = req.Email
	}
	if req.Password != "" {
		account.Password = req.Password
	}
	return account
}

func (a *AccountsHandler) fillFields(req account2.AccountRequest) *account2.Account {

	return &account2.Account{
		ID:          uuid.New().String(),
		Username:    req.Username,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		Password:    req.Password,
		CreatedAt:   time.Now().UTC().Format("2006-01-02"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02"),
		Deleted:     false,
	}
}
