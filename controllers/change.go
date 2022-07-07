package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"social_network_project/database/repository"
)

type Change struct {
	AccountRepository repository.AccountRepository
}

func (c *Change) ChangeAccount(w http.ResponseWriter, r *http.Request) {
	id, err := decodeTokenAndReturnID(r.Header.Get("BearerToken"))
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Token Invalid"}`)
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

	if err = c.AccountRepository.ChangeAccountDataByID(id, mapBody); err != nil {
		log.Fatal(err)
	}

	account, err := c.AccountRepository.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatal(err)
	}
}
