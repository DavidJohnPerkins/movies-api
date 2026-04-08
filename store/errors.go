package store

import (
	"fmt"
)

type DuplicateKeyError struct {
	ID error
}

type RecordNotFoundError struct{}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("Duplicate movie id: %v", e.ID)
}

func (e *RecordNotFoundError) Error() string {
	return "Record not found"
}
