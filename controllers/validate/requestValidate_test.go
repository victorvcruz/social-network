package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
)

func TestRequestValidate(t *testing.T) {
	validate := validator.New()

	var account = &entities.Account{
		ID:          uuid.New().String(),
		Username:    "marcelito00000000",
		Name:        "Marcelo Sabido Jose Silva",
		Description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et mas",
		Email:       "marcelo111@gmailcom",
		Password:    "23042L",
	}
	err := validate.Struct(account)
	var expectedListString1 []string
	expectedListString1 = append(expectedListString1, "Long username", "Long name", "Long description", "Invalid email", "Password only lowercase")

	listString1 := RequestValidate(err)

	assert.Equal(t, expectedListString1, listString1)

	account.Username = "marC"
	account.Name = "Ma"
	account.Description = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et ma"
	account.Email = "marcelo111@gmail.com"
	account.Password = "23042"

	err = validate.Struct(account)
	var expectedListString2 []string
	expectedListString2 = append(expectedListString2, "Username only lowercase", "Short name", "Short password")

	listString2 := RequestValidate(err)

	assert.Equal(t, expectedListString2, listString2)

	account.Username = "ma"
	account.Name = "Marcelo"
	account.Password = "230423434343443434"

	err = validate.Struct(account)
	var expectedListString3 []string
	expectedListString3 = append(expectedListString3, "Short username", "Long password")

	listString3 := RequestValidate(err)

	assert.Equal(t, expectedListString3, listString3)
}
