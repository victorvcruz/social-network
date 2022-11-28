package handlers

import (
	"github.com/stretchr/testify/assert"
	"social_network_project/internal/account"
	"testing"
	"time"
)

func TestAccountsAPI_mergeAccountToUpdatedAccount(t *testing.T) {
	var u = &account.Account{
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

	var accountExpected = &account.Account{
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

	accountExpected := &account.Account{
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
