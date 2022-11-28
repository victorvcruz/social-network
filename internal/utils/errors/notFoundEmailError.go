package errors

import "fmt"

type NotFoundEmailError struct {
	Path string
}

func (e *NotFoundEmailError) Error() string {
	return fmt.Sprintf("Incorrect email" + e.Path)
}
