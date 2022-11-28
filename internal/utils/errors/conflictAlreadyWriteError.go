package errors

import "fmt"

type ConflictAlreadyWriteError struct {
	Path string
}

func (e *ConflictAlreadyWriteError) Error() string {
	return fmt.Sprintf("Already interacted" + e.Path)
}
