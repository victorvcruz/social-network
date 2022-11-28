package crypto

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	password := "2233445566"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.Nil(t, err)
	hashExpected := string(hashedPassword)

	hashCrypt, err := EncryptPassword(password)
	assert.Nil(t, err)

	assert.Equal(t, bcrypt.CompareHashAndPassword([]byte(hashExpected), []byte(password)), bcrypt.CompareHashAndPassword([]byte(*hashCrypt), []byte(password)))
}

func TestCompareHashAndPassword(t *testing.T) {
	hash := "$2a$10$u5k/yllGsxo8rnJrkrfpLesjR1V4LoN/TRp7xYWTlZPsglw8vqQga"
	password := "2233445566"
	
	match := CompareHashAndPassword(hash, password)

	assert.Equal(t, true, match)
}
