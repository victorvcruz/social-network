package errors

import "fmt"

type ConflictEmailError struct {
	Path string
}

func (e *ConflictEmailError) Error() string {
	return fmt.Sprintf("Email already exists" + e.Path)
}
