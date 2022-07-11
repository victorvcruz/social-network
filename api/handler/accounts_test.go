package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"social_network_project/entities"
	"social_network_project/entities/response"
	"strings"
	"testing"
	"time"
)

func TestCreate_readBodyAndReturnMapBody(t *testing.T) {

	body := `{"Message": "Hello World"}`
	stringReadCloser := io.NopCloser(strings.NewReader(body))

	mapBody, err := readBodyAndReturnMapBody(stringReadCloser)
	assert.Nil(t, err)

	var mapBodyExpected map[string]interface{}

	err = json.Unmarshal([]byte(body), &mapBodyExpected)
	assert.Nil(t, err)

	assert.Equal(t, mapBodyExpected, mapBody)

}

func TestAccounts_DecodeTokenAndReturnID(t *testing.T) {
	idString := "6c08496b-b721-4e06-b0b7-1905524c9da2"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  idString,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
	assert.Nil(t, err)

	tokenStructExpected := response.Token{
		Token: tokenString,
	}

	id, err := decodeTokenAndReturnID(tokenStructExpected.Token)
	assert.Nil(t, err)

	tokenDecodeExpected := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenStructExpected.Token, tokenDecodeExpected, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
	})
	assert.Nil(t, err)

	idExpected := tokenDecodeExpected["id"].(string)

	idTest := *id

	assert.Equal(t, idExpected, idTest)
}

func TestAccountsAPI_mergeAccountToUpdatedAccount(t *testing.T) {
	var u = &entities.Account{
		ID:          "e24e073d-59e4-4fb9-9437-af65fd53f405",
		Username:    "marcelito001",
		Name:        "Marcelo Sabido",
		Description: "Eu Marcelo, Eu Marcelo",
		Email:       "marcelo111@gmail.com",
		Password:    "23042",
		CreatedAt:   "2022-07-08",
		UpdatedAt:   "2022-07-08",
		Deleted:     false,
	}

	var v = map[string]interface{}{
		"username":    "maciel",
		"name":        "Nicole Miguel Maciel",
		"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer porta vehicula purus bibendum pretium.",
		"password":    "1111111",
	}

	var accountExpected = &entities.Account{
		ID:          "e24e073d-59e4-4fb9-9437-af65fd53f405",
		Username:    "maciel",
		Name:        "Nicole Miguel Maciel",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer porta vehicula purus bibendum pretium.",
		Email:       "marcelo111@gmail.com",
		Password:    "1111111",
		CreatedAt:   "2022-07-08",
		UpdatedAt:   "2022-07-08",
		Deleted:     false,
	}

	accountUpdated := mergeAccountToUpdatedAccount(u, v)

	assert.Equal(t, accountExpected, accountUpdated)

}

func TestAccountsAPI_CreateAccount(t *testing.T) {
	var body = map[string]interface{}{
		"username":    "maciel",
		"name":        "Nicole Miguel Maciel ",
		"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer porta vehicula purus bibendum pretium.",
		"email":       "ralph333@gmail.com",
		"password":    "2222",
	}

	accountExpected := &entities.Account{
		ID:          "e1d0f3c5-3af4-4b1c-a847-1c5d8e98b2a0",
		Username:    body["username"].(string),
		Name:        body["name"].(string),
		Description: body["description"].(string),
		Email:       body["email"].(string),
		Password:    body["password"].(string),
		CreatedAt:   time.Now().UTC().Format("2006-01-02"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02"),
		Deleted:     false,
	}

	account := CreateAccountStruct(body)
	account.ID = "e1d0f3c5-3af4-4b1c-a847-1c5d8e98b2a0"

	assert.Equal(t, accountExpected, account)

}
