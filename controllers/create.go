package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"social_network_project/database/postgresql/repository"
	entities "social_network_project/entities"
	"time"
)

type Create struct {
	AccountRepository repository.AccountRepository
}

func (c *Create) CreateAccount(w http.ResponseWriter, r *http.Request) {

	mapBody, err := readBodyAndReturnMapBody(r)
	if err != nil {
		log.Fatal(err)
	}

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

	c.AccountRepository.InsertAccount(account)

	json.NewEncoder(w).Encode(account)
}

func (c *Create) CreateToken(w http.ResponseWriter, r *http.Request) {

	mapBody, err := readBodyAndReturnMapBody(r)
	if err != nil {
		log.Fatal(err)
	}

	email := mapBody["email"].(string)

	if !c.AccountRepository.ExistsAccountByEmailAndPassword(email, mapBody["password"].(string)) {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Incorrect email or password"}`)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  c.AccountRepository.FindAccountIDbyEmail(email),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte("key"))

	tokenStruct := entities.Token{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(tokenStruct)
}

func readBodyAndReturnMapBody(r *http.Request) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var mapBody map[string]interface{}

	if err := json.Unmarshal(body, &mapBody); err != nil {
		return nil, err
	}

	return mapBody, nil
}
