package request

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"social_network_project/entities"
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

	tokenString, err := token.SignedString([]byte("key"))
	assert.Nil(t, err)

	tokenStructExpected := entities.Token{
		Token: tokenString,
	}

	id, err := decodeTokenAndReturnID(tokenStructExpected.Token)
	assert.Nil(t, err)

	tokenDecodeExpected := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenStructExpected.Token, tokenDecodeExpected, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	assert.Nil(t, err)

	idExpected := tokenDecodeExpected["id"].(string)

	idTest := *id

	assert.Equal(t, idExpected, idTest)
}
