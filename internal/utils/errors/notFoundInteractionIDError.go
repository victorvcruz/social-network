package errors

import "fmt"

type NotFoundInteractionIDError struct {
	Path string
}

func (e *NotFoundInteractionIDError) Error() string {
	return fmt.Sprintf("Interaction ID does not exist" + e.Path)
}
