package errors

import "fmt"

type CacheNotFoundError struct {
	Path string
}

func (e *CacheNotFoundError) Error() string {
	return fmt.Sprintf("Not found in cache" + e.Path)
}
