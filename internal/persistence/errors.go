package persistence

import "fmt"

type ErrDuplicateKey struct {
	key string
}

func (e *ErrDuplicateKey) Error() string {
	return fmt.Sprintf("the provided %s is already being used", e.key)
}

type ErrEntityNotFound struct {
	entity string
}

func (e *ErrEntityNotFound) Error() string {
	return fmt.Sprintf("the provided %s was not found", e.entity)
}
