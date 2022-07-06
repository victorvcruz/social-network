package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"social_network_project/database/repository"
)

type Read struct {
	AccountRepository repository.AccountRepository
}

func (v *Read) GetAccount(w http.ResponseWriter, r *http.Request) {

	id, err := decodeTokenAndReturnID(r.Header.Get("BearerToken"))
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Token Invalid"}`)
		return
	}

	account, err := v.AccountRepository.FindAccountbyID(*id)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatal(err)
	}

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
