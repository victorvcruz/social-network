package request

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"social_network_project/controllers"
	"social_network_project/controllers/validate"
	"social_network_project/database/repository"
)

type AccountsRequest struct {
	AccountRepository repository.AccountRepository
	Validate          *validator.Validate
	Create            controllers.Create
}

func (a *AccountsRequest) CreateAccount(c *gin.Context) {

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var b *string
	username := mapBody["username"].(string)
	b = &username
	exist, err := a.AccountRepository.ExistsAccountByUsername(b)
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
	exist, err = a.AccountRepository.ExistsAccountByEmail(b)
	if err != nil {
		log.Fatal(err)
	}
	if *exist {
		c.JSON(http.StatusConflict, gin.H{
			"Message": "Email already exists",
		})
		return
	}

	account, err := a.Create.Account(mapBody)
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

	err = a.AccountRepository.InsertAccount(account)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account)
	return

}

func (a *AccountsRequest) CreateToken(c *gin.Context) {

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	email := mapBody["email"].(string)

	exist, err := a.AccountRepository.ExistsAccountByEmailAndPassword(email, mapBody["password"].(string))
	if err != nil {
		log.Fatal(err)
	}

	if !*exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Incorrect email or password",
		})
		return
	}

	id, err := a.AccountRepository.FindAccountIDbyEmail(email)
	if err != nil {
		log.Fatal(err)
	}

	token, err := a.Create.Token(*id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, token)
	return
}

func (a *AccountsRequest) GetAccount(c *gin.Context) {

	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountRepository.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	account, err := a.AccountRepository.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account)
	return

}

func (a *AccountsRequest) ChangeAccount(c *gin.Context) {

	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountRepository.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var b *string
	username := mapBody["username"].(string)
	b = &username
	exist, err := a.AccountRepository.ExistsAccountByUsername(b)
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
	exist, err = a.AccountRepository.ExistsAccountByEmail(b)
	if err != nil {
		log.Fatal(err)
	}
	if *exist {
		c.JSON(http.StatusConflict, gin.H{
			"Message": "Email already exists",
		})
		return
	}

	account, err := a.AccountRepository.FindAccountByID(id)
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

	if err = a.AccountRepository.ChangeAccountDataByID(id, mapBody); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, account)
	return
}

func (a *AccountsRequest) DeleteAccount(c *gin.Context) {
	id, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountRepository.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id does not exist",
		})
		return
	}

	account, err := a.AccountRepository.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	err = a.AccountRepository.DeleteAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	account.Deleted = true
	c.JSON(http.StatusOK, account)
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
		return []byte("key"), nil
	})
	if err != nil {
		return nil, err
	}
	id := tokenDecode["id"].(string)

	return &id, nil
}
