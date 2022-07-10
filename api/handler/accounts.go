package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"social_network_project/controllers"
	"social_network_project/controllers/crypto"
	"social_network_project/controllers/validate"
	"social_network_project/entities"
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
}

func (a *AccountsAPI) CreateAccount(c *gin.Context) {

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var b *string
	username := mapBody["username"].(string)
	b = &username
	exist, err := a.Controller.ExistsAccountByUsername(b)
	if err != nil {
		log.Fatal(err)
	}
	if *exist {
		c.JSON(http.StatusConflict, gin.H{
			"Message": "User already exists",
		})
		return
	}

	email := mapBody["email"].(string)
	b = &email
	exist, err = a.Controller.ExistsAccountByEmail(b)
	if err != nil {
		log.Fatal(err)
	}
	if *exist {
		c.JSON(http.StatusConflict, gin.H{
			"Message": "Email already exists",
		})
		return
	}

	account := CreateAccountStruct(mapBody)
	if err != nil {
		log.Fatal(err)
	}

	mapper := make(map[string]interface{})

	err = a.Validate.Struct(account)
	if err != nil {
		mapper["errors"] = validate.RequestValidate(err)
		c.JSON(http.StatusBadRequest, mapper)

		return
	}

	hashedPassword, err := crypto.EncryptPassword(mapBody["password"].(string))
	if err != nil {
		log.Fatal(err)
	}

	account.Password = *hashedPassword

	err = a.Controller.InsertAccount(account)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return

}

func (a *AccountsAPI) CreateToken(c *gin.Context) {

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	email := mapBody["email"].(string)

	exist, err := a.Controller.ExistsAccountByEmail(&email)
	if err != nil {
		log.Fatal(err)
	}

	if !*exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Incorrect email",
		})
		return
	}

	passwordHash, err := a.Controller.FindAccountPasswordByEmail(email)
	if err != nil {
		log.Fatal(err)
	}

	if !crypto.CompareHashAndPassword(*passwordHash, mapBody["password"].(string)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Incorrect password",
		})
		return
	}

	id, err := a.Controller.FindAccountIDbyEmail(email)
	if err != nil {
		log.Fatal(err)
	}

	token, err := CreateToken(*id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, token)
	return
}

func (a *AccountsAPI) GetAccount(c *gin.Context) {

	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.Controller.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	account, err := a.Controller.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return

}

func (a *AccountsAPI) UpdateAccount(c *gin.Context) {

	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	if mapBody["username"] != nil {
		username := mapBody["username"].(string)
		exist, err := a.Controller.ExistsAccountByUsername(&username)
		if err != nil {
			log.Fatal(err)
		}
		if *exist {
			c.JSON(http.StatusConflict, gin.H{
				"Message": "User already exists",
			})
			return
		}
	}

	if mapBody["email"] != nil {
		email := mapBody["email"].(string)
		exist, err := a.Controller.ExistsAccountByEmail(&email)
		if err != nil {
			log.Fatal(err)
		}
		if *exist {
			c.JSON(http.StatusConflict, gin.H{
				"Message": "Email already exists",
			})
			return
		}
	}

	account, err := a.Controller.FindAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	accountChange := mergeAccountToUpdatedAccount(account, mapBody)

	mapper := make(map[string]interface{})

	err = a.Validate.Struct(accountChange)
	if err != nil {
		mapper["errors"] = validate.RequestValidate(err)
		c.JSON(http.StatusBadRequest, mapper)

		return
	}

	if mapBody["password"] != nil {
		hashedPassword, err := crypto.EncryptPassword(mapBody["password"].(string))
		if err != nil {
			log.Fatal(err)
		}

		mapBody["password"] = *hashedPassword
	}

	if err = a.Controller.ChangeAccountDataByID(id, mapBody); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func (a *AccountsAPI) DeleteAccount(c *gin.Context) {
	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.Controller.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	account, err := a.Controller.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Controller.DeleteAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account.ToResponse())
	return
}

func readBodyAndReturnMapBody(body io.ReadCloser) (map[string]interface{}, error) {

	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	var mapBody map[string]interface{}

	if err := json.Unmarshal(bodyByte, &mapBody); err != nil {
		return nil, err
	}

	return mapBody, nil
}

func decodeTokenAndReturnID(token string) (*string, error) {

	tokenDecode := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, tokenDecode, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	id := tokenDecode["id"].(string)

	return &id, nil
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

	account := &entities.Account{
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

	return account
}

func CreateToken(id string) (*entities.Token, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
	if err != nil {
		return nil, err
	}

	return &entities.Token{
		Token: tokenString,
	}, nil
}
