package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social_network_project/database/repository"
)

type Delete struct {
	AccountRepository repository.AccountRepository
}

func (d *Delete) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := decodeTokenAndReturnID(r.Header.Get("BearerToken"))
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Token Invalid"}`)
		return
	}

	existID, err := d.AccountRepository.ExistsAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		w.WriteHeader(401)
		fmt.Fprintf(w, `{"Message": "Id does not exist"}`)
		return
	}

	account, err := d.AccountRepository.FindAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	err = d.AccountRepository.DeleteAccountByID(id)
	if err != nil {
		log.Fatal(err)
	}

	account.Deleted = true
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatal(err)
	}
}
