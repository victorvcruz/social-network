package errors

import "fmt"

type ConflictUsernameError struct {
	Path string
}

func (e *ConflictUsernameError) Error() string {
	return fmt.Sprintf("User already exists" + e.Path)
}
