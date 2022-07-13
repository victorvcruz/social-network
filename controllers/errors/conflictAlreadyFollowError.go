package errors

import "fmt"

type ConflictAlreadyFollowError struct {
	Path string
}

func (e *ConflictAlreadyFollowError) Error() string {
	return fmt.Sprintf("Already follow" + e.Path)
}
