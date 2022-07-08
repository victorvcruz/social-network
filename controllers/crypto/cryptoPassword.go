package crypto

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (*string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	strCrypt := string(hashedPassword)

	return &strCrypt, nil
}

func CompareHashAndPassword(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
