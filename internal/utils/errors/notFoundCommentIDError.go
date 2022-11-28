package errors

import "fmt"

type NotFoundCommentIDError struct {
	Path string
}

func (e *NotFoundCommentIDError) Error() string {
	return fmt.Sprintf("Comment ID does not exist" + e.Path)
}
