package errors

import "fmt"

type ConflictAlreadyUnfollowError struct {
	Path string
}

func (e *ConflictAlreadyUnfollowError) Error() string {
	return fmt.Sprintf("Already unfollow" + e.Path)
}
