package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"social_network_project/utils"
	"testing"
	"time"
)

func TestRequestValidate(t *testing.T) {
	t.Run("Long username, Long name, Long description, Invalid email, Password only lowercase", func(t *testing.T) {
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
		var expectedListString []string
		expectedListString = append(expectedListString, "Long username", "Long name", "Long description", "Invalid email", "Password only lowercase")
		assert.Equal(t, expectedListString, RequestAccountValidate(err))
	})
	t.Run("Username only lowercase, Short name, Short password", func(t *testing.T) {
		validate := validator.New()

		var account = &entities.Account{
			ID:          uuid.New().String(),
			Username:    "marC",
			Name:        "Ma",
			Description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et ma",
			Email:       "marcelo111@gmail.com",
			Password:    "23042",
		}
		err := validate.Struct(account)
		var expectedListString []string
		expectedListString = append(expectedListString, "Username only lowercase", "Short name", "Short password")
		assert.Equal(t, expectedListString, RequestAccountValidate(err))
	})
	t.Run("Short username, Long password", func(t *testing.T) {
		validate := validator.New()
		var account = &entities.Account{
			ID:          uuid.New().String(),
			Username:    "ma",
			Name:        "Marcelo",
			Description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et ma",
			Email:       "marcelo111@gmail.com",
			Password:    "230423434343443434",
		}
		err := validate.Struct(account)
		var expectedListString []string
		expectedListString = append(expectedListString, "Short username", "Long password")
		assert.Equal(t, expectedListString, RequestAccountValidate(err))
	})

}

func TestRequestPostValidate(t *testing.T) {
	validate := validator.New()

	var post = &entities.Post{
		ID:        "",
		AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
		Content:   "",
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	err := validate.Struct(post)
	var expectedListString1 []string
	expectedListString1 = append(expectedListString1, "Add ID", "Add content")

	listString1 := RequestPostValidate(err)

	assert.Equal(t, expectedListString1, listString1)

}

func TestRequestCommentValidate(t *testing.T) {
	validate := validator.New()

	var comment = &entities.Comment{
		ID:        "",
		AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
		PostID:    "0d0bb472-225c-4c8a-9935-a21045c80d87",
		CommentID: utils.NewNullString("8b607c43-0190-4c8c-9746-4b527d1d2c55"),
		Content:   "",
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	err := validate.Struct(comment)
	var expectedListString1 []string
	expectedListString1 = append(expectedListString1, "Add ID", "Add content")

	listString1 := RequestCommentValidate(err)

	assert.Equal(t, expectedListString1, listString1)

}

func TestRequestInteractionValidate(t *testing.T) {
	validate := validator.New()

	var interaction = &entities.Interaction{
		ID:        "",
		AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
		PostID:    utils.NewNullString("0d0bb472-225c-4c8a-9935-a21045c80d87"),
		CommentID: utils.NewNullString("8b607c43-0190-4c8c-9746-4b527d1d2c55"),
		Type:      400,
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	err := validate.Struct(interaction)
	var expectedListString1 []string
	expectedListString1 = append(expectedListString1, "Add ID", "Incorrect type, insert LIKE or DISLIKE")

	listString1 := RequestInteractionValidate(err)

	assert.Equal(t, expectedListString1, listString1)

}
