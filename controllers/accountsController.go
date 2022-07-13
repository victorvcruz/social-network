package controllers

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"social_network_project/controllers/crypto"
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
	"time"
)

type AccountsController interface {
	InsertAccount(account *entities.Account) error
	CreateToken(email string, password string) (*response.Token, error)
	FindAccountByID(id *string) (*entities.Account, error)
	ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error
	DeleteAccountByID(id *string) (*entities.Account, error)
	CreateFollow(accountID, accountToFollow *string) (*entities.Account, error)
	FindAccountsFollowing(accountID *string) ([]interface{}, error)
	FindAccountsFollowers(accountID *string) ([]interface{}, error)
	DeleteFollow(accountID, accountToFollow *string) (*entities.Account, error)
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

	existUsername, err := s.repository.ExistsAccountByUsername(&account.Username)
	if err != nil {
		return err
	}
	if *existUsername {
		return &errors.ConflictUsernameError{}
	}

	existEmail, err := s.repository.ExistsAccountByEmail(&account.Email)
	if err != nil {
		log.Fatal(err)
	}
	if *existEmail {
		return &errors.ConflictEmailError{}
	}

	hashedPassword, err := crypto.EncryptPassword(account.Password)
	if err != nil {
		return err
	}
	account.Password = *hashedPassword

	return s.repository.InsertAccount(account)
}

func (s *AccountsControllerStruct) CreateToken(email string, password string) (*response.Token, error) {

	existEmail, err := s.repository.ExistsAccountByEmail(&email)
	if err != nil {
		return nil, err
	}
	if !*existEmail {
		return nil, &errors.NotFoundEmailError{}
	}

	passwordHash, err := s.repository.FindAccountPasswordByEmail(email)
	if err != nil {
		return nil, err
	}

	if !crypto.CompareHashAndPassword(*passwordHash, password) {
		return nil, &errors.UnauthorizedPasswordError{}
	}

	id, err := s.repository.FindAccountIDbyEmail(email)
	if err != nil {
		return nil, err
	}

	token, err := CreateTokenByID(*id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AccountsControllerStruct) FindAccountByID(id *string) (*entities.Account, error) {
	account, err := s.repository.FindAccountByID(id)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return account, nil
}

func (s *AccountsControllerStruct) ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error {

	if mapBody["username"] != nil {
		username := mapBody["username"].(string)
		exist, err := s.repository.ExistsAccountByUsername(&username)
		if err != nil {
			log.Fatal(err)
		}
		if *exist {
			return &errors.ConflictUsernameError{}
		}
	}

	if mapBody["email"] != nil {
		email := mapBody["email"].(string)
		exist, err := s.repository.ExistsAccountByEmail(&email)
		if err != nil {
			return err
		}
		if *exist {
			return &errors.ConflictEmailError{}
		}
	}

	if mapBody["password"] != nil {
		hashedPassword, err := crypto.EncryptPassword(mapBody["password"].(string))
		if err != nil {
			return err
		}
		mapBody["password"] = *hashedPassword
	}

	return s.repository.ChangeAccountDataByID(id, mapBody)
}

func (s *AccountsControllerStruct) DeleteAccountByID(id *string) (*entities.Account, error) {

	account, err := s.repository.FindAccountByID(id)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	err = s.repository.DeleteAccountByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountsControllerStruct) CreateFollow(accountID, accountToFollow *string) (*entities.Account, error) {

	accountFollow, err := s.repository.FindAccountByID(accountToFollow)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	exist, err := s.repository.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.NotFoundAccountIDError{}
	}

	exist, err = s.repository.ExistsFollowByAccountIDAndAccountFollowedID(accountID, accountToFollow)
	if err != nil {
		return nil, err
	}
	if *exist {
		return nil, &errors.ConflictAlreadyFollowError{}
	}

	err = s.repository.InsertAccountFollow(accountID, accountToFollow)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return accountFollow, nil
}

func (s *AccountsControllerStruct) FindAccountsFollowing(accountID *string) ([]interface{}, error) {

	listOfAccounts, err := s.repository.FindAccountFollowingByAccountID(accountID)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return listOfAccounts, nil
}

func (s *AccountsControllerStruct) FindAccountsFollowers(accountID *string) ([]interface{}, error) {

	listOfAccounts, err := s.repository.FindAccountFollowersByAccountID(accountID)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return listOfAccounts, nil
}

func (s *AccountsControllerStruct) DeleteFollow(accountID, accountToFollow *string) (*entities.Account, error) {

	accountFollow, err := s.repository.FindAccountByID(accountToFollow)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	exist, err := s.repository.ExistsFollowByAccountIDAndAccountFollowedID(accountID, accountToFollow)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = s.repository.DeleteAccountFollow(accountID, accountToFollow)
	if err != nil {
		return nil, &errors.ConflictAlreadyUnfollowError{}
	}
	return accountFollow, nil
}

func CreateTokenByID(id string) (*response.Token, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
	if err != nil {
		return nil, err
	}

	return &response.Token{
		Token: tokenString,
	}, nil
}
