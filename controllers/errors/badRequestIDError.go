package errors

import "fmt"

type BadRequestIDError struct {
	Path string
}

func (e *BadRequestIDError) Error() string {
	return fmt.Sprintf("Add ID" + e.Path)
}
