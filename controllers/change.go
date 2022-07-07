package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"log"
	"net/http"
	"social_network_project/controllers/validate"
	"social_network_project/database/repository"
)

type Change struct {
	AccountRepository repository.AccountRepository
	Validate          *validator.Validate
}

func (c *Change) ChangeAccount(w http.ResponseWriter, r *http.Request) {
	id, err := decodeTokenAndReturnID(r.Header.Get("BearerToken"))
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Token Invalid"}`)
		return
	}

	existID, err := c.AccountRepository.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	if !*existID {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Id does not exist"}`)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	mapBody, err := readBodyAndReturnMapBody(body)
	if err != nil {
		log.Fatal(err)
	}

	account, err := c.AccountRepository.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	mapper := make(map[string]interface{})

	err = c.Validate.Struct(account)
	if err != nil {
		w.WriteHeader(400)
		mapper["errors"] = validate.RequestValidate(err)
		err = json.NewEncoder(w).Encode(mapper)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if err = c.AccountRepository.ChangeAccountDataByID(id, mapBody); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatal(err)
	}
}
