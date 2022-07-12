package errors

import "fmt"

type UnauthorizedAccountIDError struct {
	Path string
}

func (e *UnauthorizedAccountIDError) Error() string {
	return fmt.Sprintf("Unauthorized ID" + e.Path)
}
