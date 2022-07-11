package errors

import "fmt"

type NotFoundPostIDError struct {
	Path string
}

func (e *NotFoundPostIDError) Error() string {
	return fmt.Sprintf("Post ID does not exist" + e.Path)
}
