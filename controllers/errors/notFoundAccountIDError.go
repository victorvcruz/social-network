package errors

import "fmt"

type NotFoundAccountIDError struct {
	Path string
}

func (e *NotFoundAccountIDError) Error() string {
	return fmt.Sprintf("Account ID does not exist" + e.Path)
}
