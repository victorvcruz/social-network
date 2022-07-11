package validate

import (
	"github.com/go-playground/validator/v10"
)

func RequestAccountValidate(err error) []string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {

		if err.Namespace() == "Account.Username" && err.Tag() == "lowercase" {
			errors = append(errors, "Username only lowercase")
		}
		if err.Namespace() == "Account.Username" && err.Tag() == "gte" {
			errors = append(errors, "Short username")
		}
		if err.Namespace() == "Account.Username" && err.Tag() == "lte" {
			errors = append(errors, "Long username")
		}
		if err.Namespace() == "Account.Name" && err.Tag() == "gte" {
			errors = append(errors, "Short name")
		}
		if err.Namespace() == "Account.Name" && err.Tag() == "lte" {
			errors = append(errors, "Long name")
		}
		if err.Namespace() == "Account.Description" {
			errors = append(errors, "Long description")
		}
		if err.Namespace() == "Account.Email" {
			errors = append(errors, "Invalid email")
		}
		if err.Namespace() == "Account.Password" && err.Tag() == "lowercase" {
			errors = append(errors, "Password only lowercase")
		}
		if err.Namespace() == "Account.Password" && err.Tag() == "gte" {
			errors = append(errors, "Short password")
		}
		if err.Namespace() == "Account.Password" && err.Tag() == "lte" {
			errors = append(errors, "Long password")
		}
	}

	return errors
}

func RequestPostValidate(err error) []string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {

		if err.Namespace() == "Post.ID" && err.Tag() == "required" {
			errors = append(errors, "Add ID")
		}
		if err.Namespace() == "Post.Content" && err.Tag() == "required" {
			errors = append(errors, "Add content")
		}
	}

	return errors
}

func RequestCommentValidate(err error) []string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {

		if err.Namespace() == "Comment.ID" && err.Tag() == "required" {
			errors = append(errors, "Add ID")
		}
		if err.Namespace() == "Comment.Content" && err.Tag() == "required" {
			errors = append(errors, "Add content")
		}
	}

	return errors
}
