package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
	"time"
)

func TestRequestValidate(t *testing.T) {
	validate := validator.New()

	var account = &entities.Account{
		ID:          uuid.New().String(),
		Username:    "marcelito00000000",
		Name:        "Marcelo Sabido",
		Description: "I Marcelo, I Marcelo",
		Email:       "marcelo111@gmailcom",
		Password:    "23042",
		CreatedAt:   time.Now().UTC().Format("2006-01-02"),
		UpdatedAt:   time.Now().UTC().Format("2006-01-02"),
		Deleted:     false,
	}

	err := validate.Struct(account)
	var expectedListString []string
	expectedListString = append(expectedListString, "Long username", "Invalid email", "Short password")

	listString := RequestValidate(err)

	assert.Equal(t, expectedListString, listString)
}
