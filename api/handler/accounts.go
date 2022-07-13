package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"social_network_project/controllers"
	"social_network_project/controllers/errors"
	"social_network_project/controllers/validate"
	"social_network_project/entities"
	"social_network_project/utils"
	"time"
)

type AccountsAPI struct {
	Controller controllers.AccountsController
	Validate   *validator.Validate
}

func RegisterAccountsHandlers(handler *gin.Engine, accountsController controllers.AccountsController) {
	ac := &AccountsAPI{
		Controller: accountsController,
		Validate:   validator.New(),
	}

	handler.POST("/accounts", ac.CreateAccount)
	handler.POST("/accounts/auth", ac.CreateToken)
	handler.GET("/accounts", ac.GetAccount)
	handler.PUT("/accounts", ac.UpdateAccount)
	handler.DELETE("/accounts", ac.DeleteAccount)
	handler.POST("/accounts/follows", ac.FollowAccount)
	handler.GET("/accounts/following", ac.SearchFollowing)
	handler.GET("/accounts/follower", ac.SearchFollowers)
	handler.DELETE("/accounts/unfollow", ac.UnfollowAccount)

}

func (a *AccountsAPI) CreateAccount(c *gin.Context) {

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	account := CreateAccountStruct(mapBody)

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
				"Message": err.Error(),
			})
			return

		case *errors.ConflictEmailError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsAPI) CreateToken(c *gin.Context) {

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	email := mapBody["email"].(string)
	password := mapBody["password"].(string)

	token, err := a.Controller.CreateToken(email, password)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundEmailError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.UnauthorizedPasswordError:
			log.Println(e)
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, token)
	return
}

func (a *AccountsAPI) GetAccount(c *gin.Context) {

	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	account, err := a.Controller.FindAccountByID(id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return

}

func (a *AccountsAPI) UpdateAccount(c *gin.Context) {

	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	account, err := a.Controller.FindAccountByID(id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	accountChange := mergeAccountToUpdatedAccount(account, mapBody)

	mapper := make(map[string]interface{})

	err = a.Validate.Struct(accountChange)
	if err != nil {
		mapper["errors"] = validate.RequestAccountValidate(err)
		c.JSON(http.StatusBadRequest, mapper)

		return
	}

	err = a.Controller.ChangeAccountDataByID(id, mapBody)
	if err != nil {
		switch e := err.(type) {
		case *errors.ConflictUsernameError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return

		case *errors.ConflictEmailError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsAPI) DeleteAccount(c *gin.Context) {
	id, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	account, err := a.Controller.DeleteAccountByID(id)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsAPI) FollowAccount(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	if mapBody["id"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add ID",
		})
		return
	}

	accountToFollow := mapBody["id"].(string)

	account, err := a.Controller.CreateFollow(accountID, &accountToFollow)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.ConflictAlreadyFollowError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsAPI) SearchFollowing(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	listOfAccounts, err := a.Controller.FindAccountsFollowing(accountID)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, listOfAccounts)
	return
}

func (a *AccountsAPI) SearchFollowers(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	listOfAccounts, err := a.Controller.FindAccountsFollowers(accountID)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, listOfAccounts)
	return
}

func (a *AccountsAPI) UnfollowAccount(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	if mapBody["id"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add ID",
		})
		return
	}

	accountToFollow := mapBody["id"].(string)

	account, err := a.Controller.DeleteFollow(accountID, &accountToFollow)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.ConflictAlreadyUnfollowError:
			log.Println(e)
			c.JSON(http.StatusConflict, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func mergeAccountToUpdatedAccount(account *entities.Account, mapBody map[string]interface{}) *entities.Account {

	if mapBody["username"] != nil {
		account.Username = mapBody["username"].(string)
	}
	if mapBody["name"] != nil {
		account.Name = mapBody["name"].(string)
	}
	if mapBody["description"] != nil {
		account.Description = mapBody["description"].(string)
	}
	if mapBody["email"] != nil {
		account.Email = mapBody["email"].(string)
	}
	if mapBody["password"] != nil {
		account.Password = mapBody["password"].(string)
	} else {
		account.Password = "random1"

	}

	return account
}

func CreateAccountStruct(mapBody map[string]interface{}) *entities.Account {

	return &entities.Account{
		ID:          uuid.New().String(),
		Username:    mapBody["username"].(string),
		Name:        mapBody["name"].(string),
		Description: mapBody["description"].(string),
		Email:       mapBody["email"].(string),
		Password:    mapBody["password"].(string),
		CreatedAt:   time.Now().UTC().Format("2006-01-02"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02"),
		Deleted:     false,
	}
}
