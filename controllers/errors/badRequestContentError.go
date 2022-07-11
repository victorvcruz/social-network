package errors

import "fmt"

type BadRequestContentError struct {
	Path string
}

func (e *BadRequestContentError) Error() string {
	return fmt.Sprintf("Add content" + e.Path)
}
