package errors

import "fmt"

type UnauthorizedPasswordError struct {
	Path string
}

func (e *UnauthorizedPasswordError) Error() string {
	return fmt.Sprintf("Incorrect password" + e.Path)
}
