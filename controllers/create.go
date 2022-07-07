package controllers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"social_network_project/database/repository"
	entities "social_network_project/entities"
	"time"
)

type Create struct {
	AccountRepository repository.AccountRepository
}

func (c *Create) Account(mapBody map[string]interface{}) (*entities.Account, error) {

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

	return account, nil
}

func (c *Create) Token(id string) (*entities.Token, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte("key"))
	if err != nil {
		return nil, err
	}

	return &entities.Token{
		Token: tokenString,
	}, nil
}
