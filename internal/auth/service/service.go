package service

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"social_network_project/internal/account"
	"social_network_project/internal/auth"
	"social_network_project/internal/utils/crypto"
	"social_network_project/internal/utils/errors"
	"time"
)

type AuthServiceClient interface {
	CreateToken(email string, password string) (*auth.AuthResponse, error)
}

type AuthService struct {
	repository account.AccountRepository
}

func NewAuthService(accountsRepository account.AccountRepository) AuthServiceClient {
	return &AuthService{
		repository:    accountsRepository,
	}
}

func (s *AuthService) CreateToken(email string, password string) (*auth.AuthResponse, error) {

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

	token, err := s.createTokenByID(*id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) createTokenByID(id string) (*auth.AuthResponse, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
	if err != nil {
		return nil, err
	}

	return &auth.AuthResponse{
		Token: tokenString,
	}, nil
}