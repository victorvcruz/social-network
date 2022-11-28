package service

import (
	"log"
	"social_network_project/internal/account"
	"social_network_project/internal/notification"
	"social_network_project/internal/notification/service"
	"social_network_project/internal/utils/crypto"
	"social_network_project/internal/utils/errors"
)

type AccountsServiceClient interface {
	InsertAccount(account *account.Account) error
	FindAccountByID(id *string) (*account.Account, error)
	ChangeAccountDataByID(id *string, req account.AccountRequest) error
	DeleteAccountByID(id *string) (*account.Account, error)
	CreateFollow(accountID, accountToFollow *string) (*account.Account, error)
	FindAccountsFollowing(accountID, page *string) ([]interface{}, error)
	FindAccountsFollowers(accountID, page *string) ([]interface{}, error)
	DeleteFollow(accountID, accountToFollow *string) (*account.Account, error)
}

type AccountsService struct {
	repository    account.AccountRepository
	rabbitControl service.NotificationServiceClient
}

func NewAccountsService(accountsRepository account.AccountRepository, rabbitmq service.NotificationServiceClient) AccountsServiceClient {
	return &AccountsService{
		repository:    accountsRepository,
		rabbitControl: rabbitmq,
	}
}

func (s *AccountsService) InsertAccount(account *account.Account) error {

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

func (s *AccountsService) FindAccountByID(id *string) (*account.Account, error) {
	account, err := s.repository.FindAccountByID(id)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return account, nil
}

func (s *AccountsService) ChangeAccountDataByID(id *string, req account.AccountRequest) error {

	if req.Username != "" {
		username := req.Username
		exist, err := s.repository.ExistsAccountByUsername(&username)
		if err != nil {
			log.Fatal(err)
		}
		if *exist {
			return &errors.ConflictUsernameError{}
		}
	}

	if req.Email != "" {
		email := req.Email
		exist, err := s.repository.ExistsAccountByEmail(&email)
		if err != nil {
			return err
		}
		if *exist {
			return &errors.ConflictEmailError{}
		}
	}

	if req.Password != "" {
		hashedPassword, err := crypto.EncryptPassword(req.Password)
		if err != nil {
			return err
		}
		req.Password = *hashedPassword
	}

	return s.repository.ChangeAccountDataByID(id, req)
}

func (s *AccountsService) DeleteAccountByID(id *string) (*account.Account, error) {

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

func (s *AccountsService) CreateFollow(accountID, accountToFollow *string) (*account.Account, error) {

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

	s.rabbitControl.SendMessage(notification.CreateNotificationJson("FollowAccount", *accountToFollow))
	return accountFollow, nil
}

func (s *AccountsService) FindAccountsFollowing(accountID, page *string) ([]interface{}, error) {

	listOfAccounts, err := s.repository.FindAccountFollowingByAccountID(accountID, page)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return listOfAccounts, nil
}

func (s *AccountsService) FindAccountsFollowers(accountID, page *string) ([]interface{}, error) {

	listOfAccounts, err := s.repository.FindAccountFollowersByAccountID(accountID, page)
	if err != nil {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return listOfAccounts, nil
}

func (s *AccountsService) DeleteFollow(accountID, accountToFollow *string) (*account.Account, error) {

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
