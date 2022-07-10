package controllers

import (
	"social_network_project/database/repository"
	"social_network_project/entities"
)

type AccountsController interface {
	InsertAccount(account *entities.Account) error
	FindAccountPasswordByEmail(email string) (*string, error)
	FindAccountIDbyEmail(email string) (*string, error)
	FindAccountByID(id *string) (*entities.Account, error)
	ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error
	DeleteAccountByID(id *string) error
	ExistsAccountByID(id *string) (*bool, error)
	ExistsAccountByUsername(username *string) (*bool, error)
	ExistsAccountByEmail(email *string) (*bool, error)
}

type AccountsControllerStruct struct {
	repository repository.AccountRepository
}

func NewAccountsController() AccountsController {
	return &AccountsControllerStruct{
		repository: repository.NewAccountRepository(),
	}
}

func (s *AccountsControllerStruct) InsertAccount(account *entities.Account) error {
	return s.repository.InsertAccount(account)
}

func (s *AccountsControllerStruct) FindAccountPasswordByEmail(email string) (*string, error) {
	return s.repository.FindAccountPasswordByEmail(email)
}

func (s *AccountsControllerStruct) FindAccountIDbyEmail(email string) (*string, error) {
	return s.repository.FindAccountIDbyEmail(email)
}

func (s *AccountsControllerStruct) FindAccountByID(id *string) (*entities.Account, error) {
	return s.repository.FindAccountByID(id)
}

func (s *AccountsControllerStruct) ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error {
	return s.repository.ChangeAccountDataByID(id, mapBody)
}

func (s *AccountsControllerStruct) DeleteAccountByID(id *string) error {
	return s.repository.DeleteAccountByID(id)
}

func (s *AccountsControllerStruct) ExistsAccountByID(id *string) (*bool, error) {
	return s.repository.ExistsAccountByID(id)
}

func (s *AccountsControllerStruct) ExistsAccountByUsername(username *string) (*bool, error) {
	return s.repository.ExistsAccountByUsername(username)
}

func (s *AccountsControllerStruct) ExistsAccountByEmail(email *string) (*bool, error) {
	return s.repository.ExistsAccountByEmail(email)
}
